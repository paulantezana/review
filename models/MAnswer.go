package models

import "time"

// Answer struct
type Answer struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	State     bool      `json:"state" gorm:"default:'true'"`

	StudentID uint `json:"student_id"`
	PollID    uint `json:"poll_id"`

	AnswerDetails []AnswerDetail `json:"answer_details, omitempty"`
}
