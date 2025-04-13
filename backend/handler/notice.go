package handler

import (
	"net/url"
	"time"

	"github.com/kingwrcy/moments/db"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

type NoticeHandler struct {
	base BaseHandler
}

func NewNoticeHandler(injector do.Injector) *NoticeHandler {
	return &NoticeHandler{do.MustInvoke[BaseHandler](injector)}
}

// 添加公告
// @Router /api/notice/add [post]
func (n NoticeHandler) AddNotice(c echo.Context) error {
	var notice db.Notice
	if err := c.Bind(&notice); err != nil {
		return FailResp(c, ParamError)
	}

	if notice.NoticeUrl != "" {
		parsedUrl, err := url.Parse(notice.NoticeUrl)
		if err != nil || (parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https") {
			return FailRespWithMsg(c, Fail, "公告链接必须以 http 或 https 开头")
		}
	}

	now := time.Now()
	notice.CreatedAt = &now
	notice.UpdatedAt = &now

	if err := n.base.db.Create(&notice).Error; err != nil {
		return FailRespWithMsg(c, Fail, "添加公告失败")
	}

	return SuccessResp(c, notice)
}

// 获取公告列表
// @Router /api/notice/list [post]
func (n NoticeHandler) GetNoticeList(c echo.Context) error {
	var notices []db.Notice
	if err := n.base.db.Find(&notices).Error; err != nil {
		return FailRespWithMsg(c, Fail, "获取公告列表失败")
	}
	return SuccessResp(c, notices)
}
