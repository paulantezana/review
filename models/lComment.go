package models

import (
	"time"
)

type Comment struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ParentID  uint      `json:"parent_id"` // Parent comment ID
	Votes     uint32    `json:"votes"`
	HasVote   int8      `json:"has_vote" gorm:"-"` // if current user has vote
	Body      string    `json:"body"`

	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`

	User     []User    `json:"user, omitempty"`
	Children []Comment `json:"children, omitempty"`
}
