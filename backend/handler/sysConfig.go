package handler

import (
	"encoding/json"
	"errors"

	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type SysConfigHandler struct {
	base BaseHandler
}

func NewSysConfigHandler(injector do.Injector) *SysConfigHandler {
	return &SysConfigHandler{do.MustInvoke[BaseHandler](injector)}
}

// GetConfig godoc
//
//	@Tags			SysConfig
//	@Summary		获取系统设置(部分不敏感的)
//	@Description	敏感信息不返回,包括各种key密钥
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	vo.SysConfigVO
//	@Router			/sysConfig/get [post]
func (s SysConfigHandler) GetConfig(c echo.Context) error {
	var (
		config db.SysConfig
		result vo.SysConfigVO
	)

	if err := s.base.db.First(&config).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return SuccessResp(c, h{})
	}
	err := json.Unmarshal([]byte(config.Content), &result)
	if err != nil {
		return FailRespWithMsg(c, Fail, "读取系统配置异常")
	}
	result.Version = s.base.cfg.Version
	result.CommitId = s.base.cfg.CommitId

	suffix := result.S3.ThumbnailSuffix
	result.S3 = vo.S3VO{
		ThumbnailSuffix: suffix,
	}
	// 过滤敏感的OIDC信息，只保留必要字段
	result.OidcClientId = ""
	result.OidcRedirectUri = ""
	return SuccessResp(c, result)
}

// GetFullConfig godoc
//
//	@Tags		SysConfig
//	@Summary	获取系统设置(完整的)
//	@Accept		json
//	@Produce	json
//	@Param		x-api-token	header		string	true	"登录TOKEN"
//	@Success	200			{object}	vo.FullSysConfigVO
//	@Success	200
//	@Router		/api/sysConfig/get [post]
func (s SysConfigHandler) GetFullConfig(c echo.Context) error {
	var (
		config db.SysConfig
		result vo.FullSysConfigVO
	)

	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil || currentUser.Id != 1 {
		return FailRespWithMsg(c, Fail, "需要先登录")
	}
	if err := s.base.db.First(&config).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return SuccessResp(c, h{})
	}
	err := json.Unmarshal([]byte(config.Content), &result)
	if err != nil {
		return FailRespWithMsg(c, Fail, "读取系统配置异常")
	}
	result.Version = s.base.cfg.Version
	result.CommitId = s.base.cfg.CommitId
	return SuccessResp(c, result)
}

// SaveConfig godoc
//
//	@Tags		SysConfig
//	@Summary	保存系统设置
//	@Accept		json
//	@Produce	json
//	@Param		object		body	vo.FullSysConfigVO	true	"保存系统设置"
//	@Param		x-api-token	header	string				true	"登录TOKEN"
//	@Success	200
//	@Router		/api/sysConfig/save [post]
func (s SysConfigHandler) SaveConfig(c echo.Context) error {
	var (
		config db.SysConfig
		result vo.FullSysConfigVO
	)
	context := c.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil || currentUser.Id != 1 {
		s.base.log.Info().Msgf("配置保存：用户未登录或不是管理员")
		return FailRespWithMsg(c, Fail, "需要先登录")
	}

	// 记录开始保存配置
	s.base.log.Info().Msgf("开始保存系统配置，管理员ID：%d", currentUser.Id)

	if err := c.Bind(&result); err != nil {
		s.base.log.Info().Msgf("保存配置错误, %s", err)
		return FailResp(c, ParamError)
	}

	// 记录绑定的配置数据（不记录敏感信息）
	s.base.log.Info().Msgf("配置数据绑定成功，标题：%s，是否启用OIDC：%v", result.Title, result.EnableOIDC)

	data, err := json.Marshal(result)
	if err != nil {
		s.base.log.Info().Msgf("配置数据序列化失败：%s", err)
		return FailRespWithMsg(c, Fail, "读取系统配置异常")
	}

	s.base.log.Info().Msgf("配置数据序列化成功，准备保存到数据库")

	if err := s.base.db.First(&config).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		s.base.log.Info().Msgf("配置记录不存在，创建新配置记录")
		config.Content = string(data)
		if err = s.base.db.Save(&config).Error; err != nil {
			s.base.log.Info().Msgf("保存新配置记录失败：%s", err)
			return FailRespWithMsg(c, Fail, "保存系统配置异常")
		}
		s.base.log.Info().Msgf("新配置记录保存成功")
	} else {
		s.base.log.Info().Msgf("配置记录已存在，更新现有记录")
		config.Content = string(data)
		if err = s.base.db.Updates(&config).Error; err != nil {
			s.base.log.Info().Msgf("更新配置记录失败：%s", err)
			return FailRespWithMsg(c, Fail, "保存系统配置异常")
		}
		s.base.log.Info().Msgf("配置记录更新成功")
	}

	// 更新管理员用户名
	s.base.log.Info().Msgf("准备更新管理员用户名：%s", result.AdminUserName)
	s.base.db.Table("User").Where("id=?", 1).Update("username", result.AdminUserName)
	s.base.log.Info().Msgf("管理员用户名更新成功")

	return SuccessResp(c, h{})
}
