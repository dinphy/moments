package vo

import (
	"time"
)

type MessageVO struct {
	Id             int32      `json:"id,omitempty"`
	Type           string     `json:"type,omitempty"`
	Content        string     `json:"content,omitempty"`
	RelatedId      int32      `json:"relatedId,omitempty"`
	MemoId         int32      `json:"memoId,omitempty"`
	IsRead         bool       `json:"isRead,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	FromUserId     int32      `json:"fromUserId,omitempty"`
	FromName       string     `json:"fromName,omitempty"`
	FromUserAvatar string     `json:"fromUserAvatar,omitempty"`
	ReplyTo        string     `json:"replyTo,omitempty"`
}

type MessageListResp struct {
	Total int         `json:"total,omitempty"`
	List  []MessageVO `json:"list,omitempty"`
}
