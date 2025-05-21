package db

import (
	"time"
)

// Image 上传至本地的图片信息表
type Image struct {
	Id        int32      `gorm:"column:id;primary_key;NOT NULL" json:"id,omitempty"`
	CreatedAt *time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP;NOT NULL" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `gorm:"column:updatedAt;NOT NULL" json:"updatedAt,omitempty"`
	Path      string     `gorm:"column:path;NOT NULL;index:idx_img_path" json:"path,omitempty"` // 图片路径
	Hash      string     `gorm:"column:hash;NOT NULL;index:idx_img_hash" json:"hash,omitempty"` // 图片hash
}

func (i *Image) TableName() string {
	return "Image"
}

// ImageRel 图片关联表
type ImageRel struct {
	Id        int32      `gorm:"column:id;primary_key;NOT NULL" json:"id,omitempty"`
	CreatedAt *time.Time `gorm:"column:createdAt;default:CURRENT_TIMESTAMP;NOT NULL" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `gorm:"column:updatedAt;NOT NULL" json:"updatedAt,omitempty"`
	MemoId    int32      `gorm:"column:memoId;NOT NULL" json:"memoId,omitempty"`
	ImageId   int32      `gorm:"column:imageId;NOT NULL" json:"imageId,omitempty"`
}

func (i *ImageRel) TableName() string {
	return "ImageRel"
}
