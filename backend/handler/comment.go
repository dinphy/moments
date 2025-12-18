package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kingwrcy/moments/db"
	"github.com/kingwrcy/moments/pkg/mail"
	"github.com/kingwrcy/moments/pkg/util"
	"github.com/kingwrcy/moments/pkg/wechat"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

type CommentHandler struct {
	base BaseHandler
}

func NewCommentHandler(injector do.Injector) *CommentHandler {
	return &CommentHandler{do.MustInvoke[BaseHandler](injector)}
}

// RemoveComment godoc
//
//	@Tags		Comment
//	@Summary	删除评论
//	@Accept		json
//	@Produce	json
//	@Param		id			query	int		true	"评论ID"
//	@Param		x-api-token	header	string	true	"登录TOKEN"
//	@Success	200
//	@Router		/api/comment/remove [post]
func (c CommentHandler) RemoveComment(ctx echo.Context) error {
	context := ctx.(CustomContext)
	currentUser := context.CurrentUser()
	id, err := strconv.Atoi(ctx.QueryParam("id"))
	if err != nil {
		return FailResp(ctx, ParamError)
	}
	var (
		comment db.Comment
		memo    db.Memo
	)
	if err = c.base.db.First(&comment, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return FailResp(ctx, ParamError)
	}
	if err = c.base.db.First(&memo, comment.MemoId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return FailResp(ctx, ParamError)
	}

	if currentUser.Id != memo.UserId && currentUser.Id != 1 {
		return FailRespWithMsg(ctx, Fail, "没有权限")
	}
	if c.base.db.Delete(&comment).RowsAffected != 1 {
		return FailRespWithMsg(ctx, Fail, "删除失败")
	}

	// 只有当评论者不是动态发布者时，才尝试删除对应的消息
	var message db.Message
	var fromUserId int32
	var fromGuestId string

	// 获取评论者信息
	if comment.Author != "" {
		authorId, _ := strconv.ParseInt(comment.Author, 10, 32)
		fromUserId = int32(authorId)
	} else {
		fromGuestId = comment.GuestID
	}

	// 检查评论者是否是动态发布者
	if (fromUserId > 0 && fromUserId != memo.UserId) || (fromGuestId != "") {
		if err := c.base.db.Where("type = ? AND related_id = ? AND memo_id = ?", "comment", comment.Id, comment.MemoId).Delete(&message).Error; err != nil {
			c.base.log.Error().Err(err).Msg("删除评论消息失败")
		}
	}

	return SuccessResp(ctx, h{})
}

func checkGoogleRecaptcha(logger zerolog.Logger, sysConfigVO vo.FullSysConfigVO, token string) error {
	if sysConfigVO.EnableGoogleRecaptcha {
		if token == "" {
			return errors.New("token必填")
		}
		params := url.Values{}
		params.Set("secret", sysConfigVO.GoogleSecretKey)
		params.Set("response", token)

		response, err := http.Post("https://recaptcha.net/recaptcha/api/siteverify?"+params.Encode(), "", nil)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return errors.New("google验证服务无法正常返回")
		}
		resp, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		logger.Info().Str("Action", "评论").Msgf("google resp: %s", resp)

		var result map[string]any
		err = json.Unmarshal(resp, &result)
		if err != nil {
			return err
		}
		if success, ok := result["success"].(bool); ok {
			if success {
				if score, ok := result["score"].(float64); ok {
					if score > 0.5 {
						return nil
					}
				}
			}
		}
		return errors.New("人机校验不通过")
	}
	return nil
}

// AddComment godoc
//
//	@Tags		Comment
//	@Summary	添加评论
//	@Accept		json
//	@Produce	json
//	@Param		object	body	vo.AddCommentReq	true	"添加评论"
//	@Success	200
//
// @Router		/api/comment/add [post]
func (c CommentHandler) AddComment(ctx echo.Context) error {
	var (
		req         vo.AddCommentReq
		comment     db.Comment
		now         = time.Now()
		sysConfig   db.SysConfig
		sysConfigVO vo.FullSysConfigVO
	)
	err := ctx.Bind(&req)
	if err != nil {
		c.base.log.Error().Msgf("发表评论时参数校验失败,原因:%s", err)
		return FailResp(ctx, ParamError)
	}
	c.base.db.First(&sysConfig)
	_ = json.Unmarshal([]byte(sysConfig.Content), &sysConfigVO)

	if !sysConfigVO.EnableComment {
		return FailRespWithMsg(ctx, Fail, "评论未开启")
	}

	if err := checkGoogleRecaptcha(c.base.log, sysConfigVO, req.Token); err != nil {
		return FailRespWithMsg(ctx, Fail, err.Error())
	}
	if context, ok := ctx.(CustomContext); ok {
		currentUser := context.CurrentUser()
		if currentUser == nil {
			if req.GuestID == "" || !util.IsGuestID(req.GuestID) {
				return FailRespWithMsg(ctx, ParamError, "无效的访客ID")
			}
			comment.GuestID = req.GuestID

			// 使用统一的用户名处理逻辑
			isConflict := func(name string) bool {
				var count int64
				c.base.db.Model(&db.User{}).Where("username = ? OR nickname = ?", name, name).Count(&count)
				return count > 0
			}
			comment.Username = util.NormalizeUsername(req.Username, isConflict)
			if comment.Username == "" {
				comment.Username = req.GuestID
			}

			comment.Email = req.Email
			comment.Website = req.Website

			if req.Username != "" {
				// 获取该访客最新的评论姓名
				var latestComment db.Comment
				err := c.base.db.Where("guestId = ?", req.GuestID).
					Order("createdAt DESC").
					First(&latestComment).Error

				// 如果存在最新姓名，则使用该姓名更新所有记录
				if err == nil && latestComment.Username != "" {
					req.Username = latestComment.Username
				}

				// 更新同访客ID的所有评论
				c.base.db.Model(&db.Comment{}).
					Where("guestId = ?", req.GuestID).
					Update("username", req.Username)

				// 更新同访客ID的所有点赞
				c.base.db.Model(&db.Like{}).
					Where("guestId = ?", req.GuestID).
					Update("guestName", req.Username)
			}
		} else {
			comment.Username = currentUser.Nickname
			comment.Email = currentUser.Email
			comment.Author = fmt.Sprintf("%d", currentUser.Id)
		}
	}

	comment.Content = req.Content
	comment.CreatedAt = &now
	comment.UpdatedAt = &now
	comment.ReplyTo = req.ReplyTo
	comment.MemoId = req.MemoID

	// 如果是回复评论，需要获取被回复评论的邮箱
	if req.ReplyTo != "" {
		// 从数据库中查找被回复的评论
		var parentComment db.Comment
		// 通过用户名和memoId查找被回复的评论
		if err := c.base.db.Where("username = ? AND memoId = ?", req.ReplyTo, req.MemoID).First(&parentComment).Error; err == nil {
			// 仅在被回复评论的邮箱非空时使用，否则回退到前端传入的邮箱
			if parentComment.Email != "" {
				comment.ReplyEmail = parentComment.Email
			} else {
				comment.ReplyEmail = req.ReplyEmail
			}
		} else {
			// 如果找不到，使用前端传递的邮箱
			comment.ReplyEmail = req.ReplyEmail
		}
	} else {
		comment.ReplyEmail = req.ReplyEmail
	}

	if err = c.base.db.Save(&comment).Error; err == nil {
		// 创建消息通知
		var memo db.Memo
		if err := c.base.db.First(&memo, comment.MemoId).Error; err == nil {
			var fromUserId int32
			var fromGuestId string
			var fromName string

			if context, ok := ctx.(CustomContext); ok {
				currentUser := context.CurrentUser()
				if currentUser != nil {
					fromUserId = currentUser.Id
					fromName = currentUser.Nickname
				} else {
					fromGuestId = comment.GuestID
					fromName = comment.Username
				}
			}

			// 设置消息接收者
			if comment.ReplyTo != "" {
				var parentComment db.Comment
				if err := c.base.db.Where("username = ? AND memoId = ?", comment.ReplyTo, comment.MemoId).First(&parentComment).Error; err == nil {
					// 如果被回复的评论属于注册用户，则向该用户发送消息
					if parentComment.Author != "" {
						parentAuthorId, _ := strconv.ParseInt(parentComment.Author, 10, 32)
						// 避免给自己发消息
						if int32(parentAuthorId) != fromUserId {
							message := db.Message{
								UserId:      int32(parentAuthorId),
								Type:        "comment",
								Content:     comment.Content,
								RelatedId:   comment.Id,
								MemoId:      comment.MemoId,
								IsRead:      false,
								CreatedAt:   &now,
								FromUserId:  fromUserId,
								FromGuestId: fromGuestId,
								FromName:    fromName,
								ReplyTo:     comment.ReplyTo,
							}
							if err := c.base.db.Save(&message).Error; err != nil {
								c.base.log.Error().Err(err).Msg("保存消息失败")
							}
						}
					}
				}
			} else {
				// 非回复：只有当评论者不是动态发布者时才创建消息给动态发布者
				if (fromUserId > 0 && fromUserId != memo.UserId) || (fromGuestId != "") {
					message := db.Message{
						UserId:      memo.UserId,
						Type:        "comment",
						Content:     comment.Content,
						RelatedId:   comment.Id,
						MemoId:      comment.MemoId,
						IsRead:      false,
						CreatedAt:   &now,
						FromUserId:  fromUserId,
						FromGuestId: fromGuestId,
						FromName:    fromName,
						ReplyTo:     comment.ReplyTo,
					}
					if err := c.base.db.Save(&message).Error; err != nil {
						c.base.log.Error().Err(err).Msg("保存消息失败")
					}
				}
			}
		}

		go func() {
			frontendHost := fmt.Sprintf("%s://%s", ctx.Scheme(), ctx.Request().Host)
			if err = c.commentEmailNotification(comment, frontendHost); err != nil {
				c.base.log.Error().Msgf("邮件通知失败,原因:%s", err)
			}

			// 发送企业微信通知
			if err = c.sendWechatNotify(comment, frontendHost); err != nil {
				c.base.log.Error().Msgf("企业微信通知失败,原因:%s", err)
			}
		}()
		return SuccessResp(ctx, h{})
	}
	return FailRespWithMsg(ctx, Fail, "发表评论失败")
}

func (c CommentHandler) commentEmailNotification(comment db.Comment, host string) error {
	var (
		memo        db.Memo
		user        db.User
		sysConfig   db.SysConfig
		sysConfigVO vo.FullSysConfigVO
	)
	c.base.db.First(&memo, comment.MemoId)
	c.base.db.First(&user, memo.UserId)
	c.base.db.First(&sysConfig)
	_ = json.Unmarshal([]byte(sysConfig.Content), &sysConfigVO)

	// 未开启邮件通知
	if !sysConfigVO.EnableEmail {
		return nil
	}

	// 检查评论者ID（若已登录）
	var commenterUserId int32
	if comment.Author != "" {
		authorId, _ := strconv.ParseInt(comment.Author, 10, 32)
		commenterUserId = int32(authorId)
	}

	// 确定邮件接收者与目标用户ID
	var targetEmail string
	var targetUserId int32
	if comment.ReplyTo != "" { // 回复评论
		targetEmail = comment.ReplyEmail
		var parentComment db.Comment
		if err := c.base.db.Where("username = ? AND memoId = ?", comment.ReplyTo, comment.MemoId).First(&parentComment).Error; err == nil {
			if parentComment.Author != "" {
				pid, _ := strconv.ParseInt(parentComment.Author, 10, 32)
				targetUserId = int32(pid)
			}
		}
	} else { // 直接评论，通知动态发布者
		targetEmail = user.Email
		targetUserId = memo.UserId
	}

	// 如果邮件为空则直接返回
	if targetEmail == "" {
		return nil
	}

	// 如果发送者与接收者是同一已注册用户，则不发送（避免自发通知）
	if commenterUserId > 0 && targetUserId > 0 && commenterUserId == targetUserId {
		return nil
	}
	if targetEmail == "" {
		return nil
	}

	// 获取smtp客户端
	client, err := mail.GetSMTPClient(sysConfigVO.SmtpHost, sysConfigVO.SmtpPort, sysConfigVO.SmtpUsername, sysConfigVO.SmtpPassword)
	if err != nil {
		return err
	}
	defer client.Close()
	c.base.log.Info().Msgf("成功连接到SMTP服务器")

	// 通过模板生成邮件内容
	var poster string
	if comment.ReplyTo != "" { // 回复评论
		poster = comment.ReplyTo
	} else { // 直接评论
		poster = user.Nickname
	}
	data := mail.CommentNotificationEmailData{
		Title:     sysConfigVO.Title,
		Host:      host,
		Poster:    poster,
		Commenter: comment.Username,
		CommentAt: comment.CreatedAt,
		Content:   comment.Content,
		MemoId:    comment.MemoId,
	}
	emailbody, err := mail.GenerateCommentNotificationEmail(data)
	if err != nil {
		return err
	}

	getDomain := func(email string) string {
		index := strings.LastIndex(email, "@")
		domain := strings.ToLower(email[index+1:])
		return domain
	}

	// 附加头部字段
	from := sysConfigVO.SmtpUsername
	to := []string{targetEmail}
	subject := sysConfigVO.Title
	domain := getDomain(sysConfigVO.SmtpUsername)
	email := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Date: "+time.Now().Format(time.RFC1123Z)+"\r\n"+
			"Message-ID: <"+time.Now().Format("20060102150405")+"@%s>\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=utf-8\r\n"+
			"\r\n"+
			"%s",
		from, to, subject, domain, emailbody)

	// 发送邮件
	if err := client.SendMail(from, to, strings.NewReader(email)); err != nil {
		return err
	}

	c.base.log.Info().Msgf("成功发送邮件")
	return nil
}

func (c CommentHandler) sendWechatNotify(comment db.Comment, host string) error {
	var (
		memo        db.Memo
		sysConfig   db.SysConfig
		sysConfigVO vo.FullSysConfigVO
	)
	c.base.db.First(&memo, comment.MemoId)
	c.base.db.First(&sysConfig)
	_ = json.Unmarshal([]byte(sysConfig.Content), &sysConfigVO)

	// 未开启企业微信通知
	if !sysConfigVO.EnableWechatWebhook || sysConfigVO.WechatWebhookUrl == "" {
		return nil
	}

	// 构建完整的Webhook URL
	webhookURL := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", sysConfigVO.WechatWebhookUrl)

	// 检查评论者是否是动态发布者
	var commenterUserId int32
	if comment.Author != "" {
		authorId, _ := strconv.ParseInt(comment.Author, 10, 32)
		commenterUserId = int32(authorId)
	}

	// 如果评论者是动态发布者，则不发送企业微信通知
	if commenterUserId > 0 && commenterUserId == memo.UserId {
		return nil
	}

	// 构建通知标题和内容
	var title string
	var content string
	if comment.ReplyTo != "" {
		title = "💬 新的回复"
		content = fmt.Sprintf("> <font color=\"comment\">%s</font> **回复** <font color=\"info\">%s</font> 说:\n\n%s\n\n----------\n\n<font color=\"warning\">%s</font>\n\n[🔗 查看详情](%s/memo/%d)",
			comment.Username,
			comment.ReplyTo,
			comment.Content,
			comment.CreatedAt.Format("2006-01-02 15:04:05"),
			host,
			comment.MemoId)
	} else {
		title = "💬 新的评论"
		content = fmt.Sprintf("> <font color=\"comment\">%s</font> 说:\n\n%s\n\n----------\n\n<font color=\"warning\">%s</font>\n\n[🔗 查看详情](%s/memo/%d)",
			comment.Username,
			comment.Content,
			comment.CreatedAt.Format("2006-01-02 15:04:05"),
			host,
			comment.MemoId)
	}

	// 发送企业微信通知
	if err := wechat.SendWebhookNotification(webhookURL, title, content); err != nil {
		return err
	}

	return nil
}
