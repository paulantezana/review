package models

import "time"

// Reading struct
type PostReading struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Date      time.Time `json:"date"`

	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
}
