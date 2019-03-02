package models

import "time"

// AnswerDetail struct
type AnswerDetail struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Answer    string    `json:"answer"`
	State     bool      `json:"state" gorm:"default:'true'"`

	QuestionID uint `json:"question_id"`
	AnswerID   uint `json:"answer_id"`
}
