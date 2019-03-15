package models

import "time"

// AnswerDetail struct
type QuizAnswerDetail struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Answer    string    `json:"answer"`
	State     bool      `json:"state" gorm:"default:'true'"`

	QuizQuestionID uint `json:"quiz_question_id"`
	QuizAnswerID   uint `json:"quiz_answer_id"`
}
