package handler

import (
	"strconv"

	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type MessageHandler struct {
	base BaseHandler
}

func NewMessageHandler(injector do.Injector) *MessageHandler {
	return &MessageHandler{do.MustInvoke[BaseHandler](injector)}
}

// GetUnreadMessages 获取用户未读消息
// @Tags Message
// @Summary 获取用户未读消息
// @Accept json
// @Produce json
// @Param x-api-token header string true "登录TOKEN"
// @Success 200 {object} vo.MessageListResp
// @Router /api/message/unread [get]
func (m MessageHandler) GetUnreadMessages(ctx echo.Context) error {
	context := ctx.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil {
		return FailResp(ctx, 401)
	}

	var messages []db.Message
	result := m.base.db.Where("user_id = ? AND is_read = ? AND from_user_id != ?", currentUser.Id, false, currentUser.Id).
		Order("created_at DESC").
		Find(&messages)

	if result.Error != nil {
		m.base.log.Error().Err(result.Error).Msg("获取未读消息失败")
		return FailRespWithMsg(ctx, Fail, "获取未读消息失败")
	}

	var messageVOs []vo.MessageVO
	for _, msg := range messages {
		if msg.FromUserId > 0 {
			var fromUser db.User
			m.base.db.Where("id = ?", msg.FromUserId).First(&fromUser)
			msg.FromUserAvatar = fromUser.AvatarUrl
		}

		messageVOs = append(messageVOs, vo.MessageVO{
			Id:             msg.Id,
			Type:           msg.Type,
			Content:        msg.Content,
			RelatedId:      msg.RelatedId,
			MemoId:         msg.MemoId,
			IsRead:         msg.IsRead,
			CreatedAt:      msg.CreatedAt,
			FromUserId:     msg.FromUserId,
			FromName:       msg.FromName,
			FromUserAvatar: msg.FromUserAvatar,
			ReplyTo:        msg.ReplyTo,
		})
	}

	return SuccessResp(ctx, vo.MessageListResp{
		Total: len(messageVOs),
		List:  messageVOs,
	})
}

// MarkMessageAsRead 标记消息为已读
// @Tags Message
// @Summary 标记消息为已读
// @Accept json
// @Produce json
// @Param id query int true "消息ID"
// @Param x-api-token header string true "登录TOKEN"
// @Success 200
// @Router /api/message/read [post]
func (m MessageHandler) MarkMessageAsRead(ctx echo.Context) error {
	context := ctx.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil {
		return FailResp(ctx, 401)
	}

	id, err := strconv.Atoi(ctx.QueryParam("id"))
	if err != nil {
		return FailResp(ctx, ParamError)
	}

	var message db.Message
	result := m.base.db.Where("id = ? AND user_id = ?", id, currentUser.Id).First(&message)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return FailRespWithMsg(ctx, ParamError, "消息不存在")
		}
		m.base.log.Error().Err(result.Error).Msg("查找消息失败")
		return FailRespWithMsg(ctx, Fail, "查找消息失败")
	}

	message.IsRead = true
	result = m.base.db.Save(&message)
	if result.Error != nil {
		m.base.log.Error().Err(result.Error).Msg("更新消息状态失败")
		return FailRespWithMsg(ctx, Fail, "更新消息状态失败")
	}

	return SuccessResp(ctx, h{})
}

// MarkAllMessagesAsRead 标记所有消息为已读
// @Tags Message
// @Summary 标记所有消息为已读
// @Accept json
// @Produce json
// @Param x-api-token header string true "登录TOKEN"
// @Success 200
// @Router /api/message/read-all [post]
func (m MessageHandler) MarkAllMessagesAsRead(ctx echo.Context) error {
	context := ctx.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil {
		return FailResp(ctx, 401)
	}

	result := m.base.db.Model(&db.Message{}).
		Where("user_id = ? AND is_read = ?", currentUser.Id, false).
		Update("is_read", true)

	if result.Error != nil {
		m.base.log.Error().Err(result.Error).Msg("标记所有消息为已读失败")
		return FailRespWithMsg(ctx, Fail, "操作失败")
	}

	return SuccessResp(ctx, h{})
}

// GetAllMessages 获取用户所有消息
// @Tags Message
// @Summary 获取用户所有消息
// @Accept json
// @Produce json
// @Param x-api-token header string true "登录TOKEN"
// @Success 200 {object} vo.MessageListResp
// @Router /api/message/all [get]
func (m MessageHandler) GetAllMessages(ctx echo.Context) error {
	context := ctx.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil {
		return FailResp(ctx, 401)
	}

	var messages []db.Message
	result := m.base.db.Where("user_id = ? AND from_user_id != ?", currentUser.Id, currentUser.Id).
		Order("created_at DESC").
		Find(&messages)

	if result.Error != nil {
		m.base.log.Error().Err(result.Error).Msg("获取所有消息失败")
		return FailRespWithMsg(ctx, Fail, "获取所有消息失败")
	}

	var messageVOs []vo.MessageVO
	for _, msg := range messages {
		if msg.FromUserId > 0 {
			var fromUser db.User
			m.base.db.Where("id = ?", msg.FromUserId).First(&fromUser)
			msg.FromUserAvatar = fromUser.AvatarUrl
		}

		messageVOs = append(messageVOs, vo.MessageVO{
			Id:             msg.Id,
			Type:           msg.Type,
			Content:        msg.Content,
			RelatedId:      msg.RelatedId,
			MemoId:         msg.MemoId,
			IsRead:         msg.IsRead,
			CreatedAt:      msg.CreatedAt,
			FromUserId:     msg.FromUserId,
			FromName:       msg.FromName,
			FromUserAvatar: msg.FromUserAvatar,
			ReplyTo:        msg.ReplyTo,
		})
	}

	return SuccessResp(ctx, vo.MessageListResp{
		Total: len(messageVOs),
		List:  messageVOs,
	})
}

// DeleteMessage 删除单条消息
// @Tags Message
// @Summary 删除单条消息
// @Accept json
// @Produce json
// @Param id query int true "消息ID"
// @Param x-api-token header string true "登录TOKEN"
// @Success 200
// @Router /api/message/delete [delete]
func (m MessageHandler) DeleteMessage(ctx echo.Context) error {
	context := ctx.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil {
		return FailResp(ctx, 401)
	}

	id, err := strconv.Atoi(ctx.QueryParam("id"))
	if err != nil {
		return FailResp(ctx, ParamError)
	}

	var message db.Message
	result := m.base.db.Where("id = ? AND user_id = ?", id, currentUser.Id).First(&message)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return FailRespWithMsg(ctx, ParamError, "消息不存在")
		}
		m.base.log.Error().Err(result.Error).Msg("查找消息失败")
		return FailRespWithMsg(ctx, Fail, "查找消息失败")
	}

	result = m.base.db.Delete(&message)
	if result.Error != nil {
		m.base.log.Error().Err(result.Error).Msg("删除消息失败")
		return FailRespWithMsg(ctx, Fail, "删除消息失败")
	}

	return SuccessResp(ctx, h{})
}

// DeleteAllMessages 删除所有消息
// @Tags Message
// @Summary 删除所有消息
// @Accept json
// @Produce json
// @Param x-api-token header string true "登录TOKEN"
// @Success 200
// @Router /api/message/delete-all [delete]
func (m MessageHandler) DeleteAllMessages(ctx echo.Context) error {
	context := ctx.(CustomContext)
	currentUser := context.CurrentUser()
	if currentUser == nil {
		return FailResp(ctx, 401)
	}

	result := m.base.db.Where("user_id = ?", currentUser.Id).Delete(&db.Message{})

	if result.Error != nil {
		m.base.log.Error().Err(result.Error).Msg("删除所有消息失败")
		return FailRespWithMsg(ctx, Fail, "删除所有消息失败")
	}

	return SuccessResp(ctx, h{})
}
