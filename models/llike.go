package models

import "time"

type Like struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Stars     uint8     `json:"stars"`

	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`
}
