package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type OIDCHandler struct {
	base BaseHandler
}

func NewOIDCHandler(injector do.Injector) *OIDCHandler {
	return &OIDCHandler{do.MustInvoke[BaseHandler](injector)}
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

	// 构建授权URL
	authURL := s.buildAuthURL(sysConfig.OidcIssuer, sysConfig.OidcClientId, redirectURI, sysConfig.OidcScopes, state)

	return SuccessResp(c, map[string]interface{}{
		"enabled": true,
		"authURL": authURL,
		"state":   state,
	})
}

// buildAuthURL 构建OIDC授权URL
func (s *OIDCHandler) buildAuthURL(issuer, clientID, redirectURI, scopes, state string) string {
	// 构建OIDC授权端点URL
	authEndpoint := fmt.Sprintf("%s/authorize", issuer)

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
