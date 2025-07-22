package db

import (
	"time"
)

type Message struct {
	Id             int32      `gorm:"column:id;primary_key;NOT NULL" json:"id,omitempty"`
	UserId         int32      `gorm:"column:user_id;NOT NULL" json:"userId,omitempty"`               // 接收用户ID
	Type           string     `gorm:"column:type;NOT NULL" json:"type,omitempty"`                    // 消息类型: comment, like
	Content        string     `gorm:"column:content;NOT NULL" json:"content,omitempty"`              // 消息内容
	RelatedId      int32      `gorm:"column:related_id" json:"relatedId,omitempty"`                  // 相关ID (评论ID或点赞ID)
	MemoId         int32      `gorm:"column:memo_id" json:"memoId,omitempty"`                        // 相关动态ID
	IsRead         bool       `gorm:"column:is_read;default:false;NOT NULL" json:"isRead,omitempty"` // 是否已读
	CreatedAt      *time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"createdAt,omitempty"`
	FromUserId     int32      `gorm:"column:from_user_id" json:"fromUserId,omitempty"`         // 发送用户ID
	FromGuestId    string     `gorm:"column:from_guest_id" json:"fromGuestId,omitempty"`       // 发送访客ID
	FromName       string     `gorm:"column:from_name;NOT NULL" json:"fromName,omitempty"`     // 发送者名称
	FromUserAvatar string     `gorm:"column:from_user_avatar" json:"fromUserAvatar,omitempty"` // 发送用户头像
}

func (m *Message) TableName() string {
	return "Message"
}
