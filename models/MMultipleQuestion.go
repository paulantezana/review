package models

import "time"

// MultipleQuestion struct
type MultipleQuestion struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Label     string    `json:"label"`
	State     bool      `json:"state" gorm:"default:'true'"`

	QuestionID uint `json:"question_id"`
}
