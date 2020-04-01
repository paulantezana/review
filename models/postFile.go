package models

import "time"

type PostFile struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint      `json:"-"`
	UpdateUserId uint      `json:"-"`

	Path           string `json:"path"`
	PostID         uint   `json:"post_id"`
	EnableDownload bool   `json:"enable_download"`
}
