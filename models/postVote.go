package models

import "time"

type PostVote struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostCommentID uint `json:"post_comment_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Value     bool      `json:"value" gorm:"not null"`
}
