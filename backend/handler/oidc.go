package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/pkg/oidc"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type OIDCHandler struct {
	base    BaseHandler
	oidcSvc oidc.Service
}

func NewOIDCHandler(injector do.Injector) *OIDCHandler {
	return &OIDCHandler{
		base:    do.MustInvoke[BaseHandler](injector),
		oidcSvc: oidc.NewOIDCService(),
	}
}

// GetOIDCConfig godoc
//
// @Tags OIDC
// @Summary 获取OIDC配置信息
// @Description 获取OIDC配置信息，用于前端判断是否显示OIDC登录按钮及获取登录地址
// @Accept json
// @Produce json
// @Success 200 {object} vo.OIDCConfigResponse
// @Router /api/oidc/config [post]
func (s *OIDCHandler) GetOIDCConfig(c echo.Context) error {
	// 查询系统配置
	var config db.SysConfig
	if err := s.base.db.First(&config).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		// 配置不存在，返回OIDC未启用
		return SuccessResp(c, &vo.OIDCConfigResponse{
			Enabled: false,
		})
	}

	// 解析配置
	var sysConfig vo.SysConfigVO
	err := json.Unmarshal([]byte(config.Content), &sysConfig)
	if err != nil {
		s.base.log.Info().Msgf("解析系统配置失败: %s", err)
		return SuccessResp(c, &vo.OIDCConfigResponse{
			Enabled: false,
		})
	}

	// 检查OIDC是否启用
	if !sysConfig.EnableOIDC || sysConfig.OidcIssuer == "" || sysConfig.OidcClientId == "" {
		return SuccessResp(c, &vo.OIDCConfigResponse{
			Enabled: false,
		})
	}

	// 构建OIDC授权URL
	redirectURI := sysConfig.OidcRedirectUri
	if redirectURI == "" {
		// 如果未配置redirectURI，使用默认值
		host := c.Request().Host
		redirectURI = fmt.Sprintf("https://%s/api/oidc/callback", host)
	}

	// 使用OIDC服务获取配置
	oidcConfig, err := s.oidcSvc.GetConfig(
		sysConfig.OidcIssuer,
		sysConfig.OidcClientId,
		redirectURI,
		sysConfig.OidcScopes,
	)

	if err != nil {
		s.base.log.Warn().Msgf("获取OIDC配置失败: %v", err)
		return SuccessResp(c, &vo.OIDCConfigResponse{
			Enabled:     true,
			ConfigValid: false,
		})
	}

	return SuccessResp(c, oidcConfig)
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

	// 使用OIDC服务处理回调
	systemUser, err := s.processUserInfoWithOIDC(code, state, sysConfig.OidcIssuer, sysConfig.OidcClientId, sysConfig.OidcClientSecret, redirectURI)
	if err != nil {
		s.base.log.Error().Msgf("处理OIDC回调失败: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "OIDC回调处理失败",
		})
	}

	// 生成JWT token
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

// processUserInfoWithOIDC 处理OIDC用户信息
func (s *OIDCHandler) processUserInfoWithOIDC(code, state, issuer, clientId, clientSecret, redirectURI string) (*db.User, error) {
	// 使用OIDC服务处理回调
	userInfo, _, err := s.oidcSvc.HandleCallback(code, state, issuer, clientId, clientSecret, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("OIDC服务处理回调失败: %w", err)
	}

	// 处理用户信息，创建或更新系统用户
	systemUser, err := s.processUserInfo(userInfo)
	if err != nil {
		return nil, fmt.Errorf("处理用户信息失败: %w", err)
	}

	// 更新系统用户信息
	s.base.log.Info().Msgf("获取用户信息%s", systemUser)

	return systemUser, nil
}

// processUserInfo 处理OIDC用户信息，创建或更新系统用户
func (s *OIDCHandler) processUserInfo(userInfo *vo.OIDCUserInfo) (*db.User, error) {
	// 确定昵称，优先使用Name，如果没有则使用NickName
	nickname := userInfo.Name
	if nickname == "" {
		nickname = userInfo.NickName
	}

	// 检查Id为1的用户的OidcSub是否为空
	var adminUser db.User
	adminUserResult := s.base.db.Where("id = ?", 1).First(&adminUser)

	// 如果Id为1的用户存在且其OidcSub为空，则将当前OIDC信息绑定到该用户
	if adminUserResult.Error == nil && adminUser.OidcSub == "" {
		// 更新用户信息，绑定OIDC
		adminUser.OidcSub = userInfo.Sub
		s.base.db.Save(&adminUser)
		return &adminUser, nil
	}

	// 如果管理员账号已经绑定了OIDC，且当前登录的OIDC用户与管理员账号匹配，则更新管理员信息
	if adminUserResult.Error == nil && adminUser.OidcSub == userInfo.Sub {
		return &adminUser, nil
	}

	// 查找是否已存在该OIDC用户
	var existingUser db.User
	result := s.base.db.Where("oidcSub = ?", userInfo.Sub).First(&existingUser)

	// 如果用户已存在，返回现有用户
	if result.Error == nil {
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
		Nickname:  nickname,
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
