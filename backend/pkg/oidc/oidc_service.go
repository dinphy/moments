package oidc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/vo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Service OIDC服务接口
type Service interface {
	GetConfig(issuer, clientId, redirectURI, scopes string) (*vo.OIDCConfigResponse, error)
	GetUserInfo(code, state, issuer, clientId, clientSecret, redirectURI string) (*vo.OIDCUserInfo, error)
	HandleCallback(code, state, issuer, clientId, clientSecret, redirectURI string) (*vo.OIDCUserInfo, string, error)
}

type oidcService struct {
	configCache OIDCConfigCache
}

// NewOIDCService 创建OIDC服务实例
func NewOIDCService() Service {
	return &oidcService{
		configCache: OIDCConfigCache{},
	}
}

// GetConfig 获取OIDC配置信息
func (s *oidcService) GetConfig(issuer, clientId, redirectURI, scopes string) (*vo.OIDCConfigResponse, error) {
	// 生成随机state参数，用于后续验证
	state := uuid.New().String()

	// 尝试验证OIDC配置是否有效
	_, oidcErr := s.getOIDCMetadata(issuer)

	// 构建授权URL
	authURL := s.buildAuthURL(issuer, clientId, redirectURI, scopes, state)

	return &vo.OIDCConfigResponse{
		Enabled:     true,
		AuthURL:     authURL,
		State:       state,
		ConfigValid: oidcErr == nil,
	}, nil
}

// HandleCallback 处理OIDC回调
func (s *oidcService) HandleCallback(code, state, issuer, clientId, clientSecret, redirectURI string) (*vo.OIDCUserInfo, string, error) {
	// 获取access token
	tokenResponse, err := s.getAccessToken(issuer, clientId, clientSecret, code, redirectURI)
	if err != nil {
		return nil, "", fmt.Errorf("获取access token失败: %w", err)
	}

	// 获取用户信息
	userInfo, err := s.getUserInfo(issuer, tokenResponse.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 生成JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userInfo.Sub,
	})

	// 这里需要从外部传入jwtKey，暂时使用空字符串
	tokenString, err := token.SignedString([]byte(""))
	if err != nil {
		return nil, "", fmt.Errorf("生成JWT token失败: %w", err)
	}

	// 返回用户信息和token，但不处理数据库相关操作，由Handler处理
	return userInfo, tokenString, nil
}

// GetUserInfo 获取OIDC用户信息
func (s *oidcService) GetUserInfo(code, state, issuer, clientId, clientSecret, redirectURI string) (*vo.OIDCUserInfo, error) {
	// 获取access token
	tokenResponse, err := s.getAccessToken(issuer, clientId, clientSecret, code, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("获取access token失败: %w", err)
	}

	// 获取用户信息
	userInfo, err := s.getUserInfo(issuer, tokenResponse.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	return userInfo, nil
}

// getOIDCMetadata 获取OIDC提供商元数据
func (s *oidcService) getOIDCMetadata(issuer string) (*vo.OIDCMetadata, error) {
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

	resp, err := http.Get(configURL)
	if err != nil {
		return nil, fmt.Errorf("获取OIDC配置失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取OIDC配置失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取OIDC配置响应失败: %w", err)
	}

	var metadata vo.OIDCMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
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
func (s *oidcService) buildAuthURL(issuer, clientID, redirectURI, scopes, state string) string {
	// 尝试从OIDC配置获取授权端点
	metadata, err := s.getOIDCMetadata(issuer)
	var authEndpoint string

	if err != nil || metadata.AuthorizationEndpoint == "" {
		// 如果获取配置失败或授权端点为空，使用默认路径作为备选
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

// getAccessToken 通过授权码获取access token
func (s *oidcService) getAccessToken(issuer, clientID, clientSecret, code, redirectURI string) (*vo.OIDCTokenResponse, error) {
	// 尝试从OIDC配置获取token端点
	metadata, err := s.getOIDCMetadata(issuer)
	var tokenEndpoint string

	if err != nil || metadata.TokenEndpoint == "" {
		// 如果获取配置失败或token端点为空，使用默认路径作为备选
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
		return nil, fmt.Errorf("请求token失败，状态码: %d, 响应: %s", resp.StatusCode, string(errorBody))
	}

	// 解析响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResponse vo.OIDCTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, fmt.Errorf("解析token响应失败: %w", err)
	}

	return &tokenResponse, nil
}

// getUserInfo 通过access token获取用户信息
func (s *oidcService) getUserInfo(issuer, accessToken string) (*vo.OIDCUserInfo, error) {
	// 尝试从OIDC配置获取userinfo端点
	metadata, err := s.getOIDCMetadata(issuer)
	var userInfoEndpoint string

	if err != nil || metadata.UserinfoEndpoint == "" {
		// 如果获取配置失败或userinfo端点为空，使用默认路径作为备选
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
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		// 读取错误响应体以提供更详细的错误信息
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("请求用户信息失败，状态码: %d, 响应: %s", resp.StatusCode, string(errorBody))
	}

	// 解析响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo vo.OIDCUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	// 验证用户信息是否包含必要的sub字段
	if userInfo.Sub == "" {
		return nil, fmt.Errorf("用户信息缺少必要的sub字段")
	}

	return &userInfo, nil
}

// processUserInfo 处理OIDC用户信息，创建或更新系统用户
// 注意：此方法需要数据库访问，应在Handler中实现
// 此处仅作为接口定义，实际实现由Handler提供
func (s *oidcService) ProcessUserInfo(userInfo *vo.OIDCUserInfo, dbHandler interface{}) (*db.User, error) {
	// 这个方法需要访问数据库，应该在Handler中实现
	// 这里只返回一个错误，提示需要在Handler中实现
	return nil, fmt.Errorf("ProcessUserInfo方法需要在Handler中实现，因为需要访问数据库")
}

// OIDCConfigCache OIDC配置缓存结构
type OIDCConfigCache struct {
	metadata   *vo.OIDCMetadata
	expiration time.Time
	lastIssuer string
	lock       sync.RWMutex
}
