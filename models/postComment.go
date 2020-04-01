package models

import (
	"time"
)

type PostComment struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ParentID  uint      `json:"parent_id"` // Parent comment ID
    PostVotes     uint32 `json:"post_votes"`
	HasPostVote   int8      `json:"has_post_vote" gorm:"-"` // if current user has vote
	Body      string    `json:"body"`

	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`

	User     []User        `json:"user, omitempty"`
	Children []PostComment `json:"children, omitempty"`
}
