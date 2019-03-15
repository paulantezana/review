package models

import "time"

// Answer struct
type QuizAnswer struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	State     bool      `json:"state" gorm:"default:'true'"`

	StudentID uint `json:"student_id"`
	QuizID    uint `json:"quiz_id"`

	AnswerDetails []QuizAnswerDetail `json:"answer_details, omitempty"`
}
