package models

import "time"

// Reading struct
type Reading struct {
	ID   uint      `json:"id" gorm:"primary_key"`
	Date time.Time `json:"date"`

	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`
}
