package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type OIDCHandler struct {
	base        BaseHandler
	configCache OIDCConfigCache
}

func NewOIDCHandler(injector do.Injector) *OIDCHandler {
	return &OIDCHandler{
		base:        do.MustInvoke[BaseHandler](injector),
		configCache: OIDCConfigCache{},
	}
}

// GetOIDCConfig godoc
//
// @Tags OIDC
// @Summary 获取OIDC配置信息
// @Description 获取OIDC配置信息，用于前端判断是否显示OIDC登录按钮及获取登录地址
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/oidc/config [post]
func (s *OIDCHandler) GetOIDCConfig(c echo.Context) error {
	// 查询系统配置
	var config db.SysConfig
	if err := s.base.db.First(&config).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// 配置不存在，返回OIDC未启用
		return SuccessResp(c, map[string]interface{}{
			"enabled": false,
		})
	}

	// 解析配置
	var sysConfig vo.SysConfigVO
	err := json.Unmarshal([]byte(config.Content), &sysConfig)
	if err != nil {
		s.base.log.Info().Msgf("解析系统配置失败: %s", err)
		return SuccessResp(c, map[string]interface{}{
			"enabled": false,
		})
	}

	// 检查OIDC是否启用
	if !sysConfig.EnableOIDC || sysConfig.OidcIssuer == "" || sysConfig.OidcClientId == "" {
		return SuccessResp(c, map[string]interface{}{
			"enabled": false,
		})
	}

	// 生成随机state参数，用于后续验证
	state := uuid.New().String()

	// 构建OIDC授权URL
	redirectURI := sysConfig.OidcRedirectUri
	if redirectURI == "" {
		// 如果未配置redirectURI，使用默认值
		host := c.Request().Host
		redirectURI = fmt.Sprintf("https://%s/api/oidc/callback", host)
	}

	// 尝试验证OIDC配置是否有效
	_, oidcErr := s.getOIDCMetadata(sysConfig.OidcIssuer)
	if oidcErr != nil {
		// OIDC配置无效，但不阻止返回配置信息，仅记录警告
		s.base.log.Warn().Msgf("OIDC配置验证失败，但仍返回配置信息: %v", oidcErr)
	}

	// 构建授权URL
	authURL := s.buildAuthURL(sysConfig.OidcIssuer, sysConfig.OidcClientId, redirectURI, sysConfig.OidcScopes, state)

	return SuccessResp(c, map[string]interface{}{
		"enabled":     true,
		"authURL":     authURL,
		"state":       state,
		"configValid": oidcErr == nil, // 添加配置有效性标志
	})
}

// getOIDCMetadata 获取OIDC提供商元数据
func (s *OIDCHandler) getOIDCMetadata(issuer string) (*OIDCMetadata, error) {
	// 检查缓存
	s.configCache.lock.RLock()
	if s.configCache.metadata != nil &&
		s.configCache.lastIssuer == issuer &&
		time.Now().Before(s.configCache.expiration) {
		metadata := s.configCache.metadata
		s.configCache.lock.RUnlock()
		return metadata, nil
	}
	s.configCache.lock.RUnlock()

	// 缓存不存在或已过期，从OIDC提供商获取
	configURL := fmt.Sprintf("%s/.well-known/openid-configuration", issuer)
	s.base.log.Info().Msgf("获取OIDC配置信息: %s", configURL)

	resp, err := http.Get(configURL)
	if err != nil {
		s.base.log.Error().Msgf("获取OIDC配置失败: %s", err)
		return nil, fmt.Errorf("获取OIDC配置失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.base.log.Error().Msgf("获取OIDC配置失败，状态码: %d", resp.StatusCode)
		return nil, fmt.Errorf("获取OIDC配置失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.base.log.Error().Msgf("读取OIDC配置响应失败: %s", err)
		return nil, fmt.Errorf("读取OIDC配置响应失败: %w", err)
	}

	var metadata OIDCMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		s.base.log.Error().Msgf("解析OIDC配置失败: %s", err)
		return nil, fmt.Errorf("解析OIDC配置失败: %w", err)
	}

	// 更新缓存，有效期1小时
	s.configCache.lock.Lock()
	s.configCache.metadata = &metadata
	s.configCache.lastIssuer = issuer
	s.configCache.expiration = time.Now().Add(time.Hour)
	s.configCache.lock.Unlock()

	return &metadata, nil
}

// buildAuthURL 构建OIDC授权URL
func (s *OIDCHandler) buildAuthURL(issuer, clientID, redirectURI, scopes, state string) string {
	// 尝试从OIDC配置获取授权端点
	metadata, err := s.getOIDCMetadata(issuer)
	var authEndpoint string

	if err != nil || metadata.AuthorizationEndpoint == "" {
		// 如果获取配置失败或授权端点为空，使用默认路径作为备选
		s.base.log.Warn().Msgf("使用默认授权端点路径，原因: %v", err)
		authEndpoint = fmt.Sprintf("%s/auth", issuer)
	} else {
		authEndpoint = metadata.AuthorizationEndpoint
	}

	// 创建查询参数
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("scope", scopes)
	params.Add("state", state)

	// 返回完整的授权URL
	return fmt.Sprintf("%s?%s", authEndpoint, params.Encode())
}

// OIDCMetadata OIDC提供商元数据结构
type OIDCMetadata struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
	RevocationEndpoint                string   `json:"revocation_endpoint,omitempty"`
	IntrospectionEndpoint             string   `json:"introspection_endpoint,omitempty"`
	JwksURI                           string   `json:"jwks_uri"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	ScopesSupported                   []string `json:"scopes_supported,omitempty"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	ClaimsSupported                   []string `json:"claims_supported,omitempty"`
}

// OIDCTokenResponse OIDC token响应结构
type OIDCTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	IDToken     string `json:"id_token,omitempty"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}

// OIDCUserInfo OIDC用户信息结构
type OIDCUserInfo struct {
	Sub               string `json:"sub"`  // 用户唯一标识
	Name              string `json:"name"` // 用户名称
	NickName          string `json:"nickname,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Email             string `json:"email,omitempty"`
	Picture           string `json:"picture,omitempty"`
}

// OIDCConfigCache OIDC配置缓存结构
type OIDCConfigCache struct {
	metadata   *OIDCMetadata
	expiration time.Time
	lastIssuer string
	lock       sync.RWMutex
}

// GetOIDCCallback godoc
//
// @Tags OIDC
// @Summary OIDC登录回调处理
// @Description 处理OIDC提供商的授权回调，获取用户信息并完成登录
// @Accept json
// @Produce json
// @Param code query string true "授权码"
// @Param state query string true "状态参数"
// @Router /api/oidc/callback [get]
func (s *OIDCHandler) GetOIDCCallback(c echo.Context) error {
	// 获取查询参数
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" || state == "" {
		s.base.log.Warn().Msgf("OIDC回调缺少必要参数，code=%s, state=%s", code, state)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "缺少必要参数",
		})
	}

	// 查询系统配置
	var config db.SysConfig
	if err := s.base.db.First(&config).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		s.base.log.Error().Msgf("OIDC回调时系统配置不存在")
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "系统配置不存在",
		})
	}

	// 解析配置
	var sysConfig vo.SysConfigVO
	err := json.Unmarshal([]byte(config.Content), &sysConfig)
	if err != nil {
		s.base.log.Error().Msgf("OIDC回调时解析系统配置失败: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "系统配置解析失败",
		})
	}

	// 检查OIDC是否启用
	if !sysConfig.EnableOIDC || sysConfig.OidcIssuer == "" || sysConfig.OidcClientId == "" {
		s.base.log.Warn().Msgf("OIDC回调时OIDC未正确启用")
		return c.JSON(http.StatusNotImplemented, map[string]interface{}{
			"error": "OIDC未启用",
		})
	}

	// 构建redirectURI
	redirectURI := sysConfig.OidcRedirectUri
	if redirectURI == "" {
		host := c.Request().Host
		redirectURI = fmt.Sprintf("https://%s/api/oidc/callback", host)
	}

	// 获取access token
	tokenResponse, err := s.getAccessToken(sysConfig.OidcIssuer, sysConfig.OidcClientId, sysConfig.OidcClientSecret, code, redirectURI)
	if err != nil {
		s.base.log.Error().Msgf("获取access token失败: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "获取访问令牌失败",
		})
	}

	// 获取用户信息
	userInfo, err := s.getUserInfo(sysConfig.OidcIssuer, tokenResponse.AccessToken)
	if err != nil {
		s.base.log.Error().Msgf("获取用户信息失败: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "获取用户信息失败",
		})
	}

	// 处理用户信息，创建或更新系统用户
	systemUser, err := s.processUserInfo(userInfo, sysConfig)
	if err != nil {
		s.base.log.Error().Msgf("处理用户信息失败: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "用户信息处理失败",
		})
	}
	s.base.log.Info().Msgf("获取用户信息%s", systemUser)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": systemUser.Username,
		"userId":   systemUser.Id,
	})

	tokenString, err := token.SignedString([]byte(s.base.cfg.JwtKey))
	if err != nil {
		s.base.log.Error().Msgf("生成jwt token异常:%s", err)
		return FailRespWithMsg(c, Fail, "登录异常")
	}
	// 重定向到OIDC回调页面，并在URL参数中携带token、userId和username
	return c.Redirect(http.StatusFound, fmt.Sprintf("/oidc?token=%s&userId=%d&username=%s", tokenString, systemUser.Id, systemUser.Username))
}

// getAccessToken 通过授权码获取access token
func (s *OIDCHandler) getAccessToken(issuer, clientID, clientSecret, code, redirectURI string) (*OIDCTokenResponse, error) {
	// 尝试从OIDC配置获取token端点
	metadata, err := s.getOIDCMetadata(issuer)
	var tokenEndpoint string

	if err != nil || metadata.TokenEndpoint == "" {
		// 如果获取配置失败或token端点为空，使用默认路径作为备选
		s.base.log.Warn().Msgf("使用默认token端点路径，原因: %v", err)
		tokenEndpoint = fmt.Sprintf("%s/token", issuer)
	} else {
		tokenEndpoint = metadata.TokenEndpoint
	}

	// 准备请求体
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", redirectURI)
	data.Add("client_id", clientID)
	if clientSecret != "" {
		data.Add("client_secret", clientSecret)
	}

	// 发送POST请求
	req, err := http.NewRequest("POST", tokenEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		// 读取错误响应体以提供更详细的错误信息
		errorBody, _ := io.ReadAll(resp.Body)
		s.base.log.Error().Msgf("请求token失败，状态码: %d, 响应: %s", resp.StatusCode, string(errorBody))
		return nil, fmt.Errorf("请求token失败，状态码: %d, 响应: %s", resp.StatusCode, string(errorBody))
	}

	// 解析响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResponse OIDCTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		s.base.log.Error().Msgf("解析token响应失败: %s, 响应体: %s", err, string(body))
		return nil, fmt.Errorf("解析token响应失败: %w", err)
	}

	return &tokenResponse, nil
}

// getUserInfo 通过access token获取用户信息
func (s *OIDCHandler) getUserInfo(issuer, accessToken string) (*OIDCUserInfo, error) {
	// 尝试从OIDC配置获取userinfo端点
	metadata, err := s.getOIDCMetadata(issuer)
	var userInfoEndpoint string

	if err != nil || metadata.UserinfoEndpoint == "" {
		// 如果获取配置失败或userinfo端点为空，使用默认路径作为备选
		s.base.log.Warn().Msgf("使用默认userinfo端点路径，原因: %v", err)
		userInfoEndpoint = fmt.Sprintf("%s/userinfo", issuer)
	} else {
		userInfoEndpoint = metadata.UserinfoEndpoint
	}

	// 发送GET请求
	req, err := http.NewRequest("GET", userInfoEndpoint, nil)
	if err != nil {
		return nil, err
	}

	// 设置Authorization头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.base.log.Error().Msgf("请求用户信息失败: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		// 读取错误响应体以提供更详细的错误信息
		errorBody, _ := io.ReadAll(resp.Body)
		s.base.log.Error().Msgf("请求用户信息失败，状态码: %d, 响应: %s", resp.StatusCode, string(errorBody))
		return nil, fmt.Errorf("请求用户信息失败，状态码: %d, 响应: %s", resp.StatusCode, string(errorBody))
	}

	// 解析响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo OIDCUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		s.base.log.Error().Msgf("解析用户信息失败: %s, 响应体: %s", err, string(body))
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	// 验证用户信息是否包含必要的sub字段
	if userInfo.Sub == "" {
		s.base.log.Error().Msgf("用户信息缺少必要的sub字段: %s", string(body))
		return nil, fmt.Errorf("用户信息缺少必要的sub字段")
	}

	return &userInfo, nil
}

// processUserInfo 处理OIDC用户信息，创建或更新系统用户
func (s *OIDCHandler) processUserInfo(userInfo *OIDCUserInfo, sysConfig vo.SysConfigVO) (*db.User, error) {
	// 检查Id为1的用户的OidcSub是否为空
	var adminUser db.User
	adminUserResult := s.base.db.Where("id = ?", 1).First(&adminUser)

	// 如果Id为1的用户存在且其OidcSub为空，则将当前OIDC信息绑定到该用户
	if adminUserResult.Error == nil && adminUser.OidcSub == "" {
		// 更新用户信息，绑定OIDC
		adminUser.OidcSub = userInfo.Sub
		adminUser.Nickname = userInfo.NickName
		if userInfo.PreferredUsername != "" {
			adminUser.Username = userInfo.PreferredUsername
		}
		if userInfo.Email != "" {
			adminUser.Email = userInfo.Email
		}
		if userInfo.Picture != "" {
			adminUser.AvatarUrl = userInfo.Picture
		}
		s.base.db.Save(&adminUser)
		return &adminUser, nil
	}

	// 查找是否已存在该OIDC用户
	var existingUser db.User
	result := s.base.db.Where("oidcSub = ?", userInfo.Sub).First(&existingUser)

	// 如果用户已存在，返回现有用户
	if result.Error == nil {
		// 更新用户信息
		existingUser.Nickname = userInfo.Name
		if userInfo.PreferredUsername != "" {
			existingUser.Username = userInfo.PreferredUsername
		}
		if userInfo.Email != "" {
			existingUser.Email = userInfo.Email
		}
		if userInfo.Picture != "" {
			existingUser.AvatarUrl = userInfo.Picture
		}
		s.base.db.Save(&existingUser)
		return &existingUser, nil
	}

	// 如果用户不存在且不是errors.Is(result.Error, gorm.ErrRecordNotFound)，返回错误
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 创建新用户
	username := userInfo.PreferredUsername
	if username == "" {
		// 如果没有preferred_username，使用sub作为用户名
		username = fmt.Sprintf("oidc_%s", userInfo.Sub)
		// 确保用户名唯一
		var count int64
		s.base.db.Model(&db.User{}).Where("username = ?", username).Count(&count)
		if count > 0 {
			username = fmt.Sprintf("%s_%d", username, count+1)
		}
	}

	newUser := db.User{
		OidcSub:   userInfo.Sub,
		Username:  username,
		Nickname:  userInfo.NickName,
		Email:     userInfo.Email,
		AvatarUrl: userInfo.Picture,
		Password:  "", // OIDC用户不需要密码
		CoverUrl:  "/cover.webp",
		Favicon:   "/favicon.png",
	}

	// 保存用户
	if err := s.base.db.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}
