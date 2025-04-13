package db

import (
    "time"
)

type Notice struct {
    Id          int32  `gorm:"column:id;primary_key;NOT NULL" json:"id,omitempty"`
    Title       string `gorm:"column:title" json:"title,omitempty"`
    Content     string `gorm:"column:content" json:"content,omitempty"`
    NoticeUrl   string `gorm:"column:noticeUrl" json:"noticeUrl,omitempty"`
    Description string `gorm:"column:description" json:"description,omitempty"`
    CreatedAt   *time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP;NOT NULL" json:"createdAt,omitempty"`
    UpdatedAt   *time.Time `gorm:"column:updatedAt;NOT NULL" json:"updatedAt,omitempty"`
}

func (n *Notice) TableName() string {
    return "Notice"
}