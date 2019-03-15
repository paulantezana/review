package models

import "time"

// MultipleQuestion struct
type MultipleQuizQuestion struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Label     string    `json:"label"`
	Correct   bool      `json:"correct"`
	State     bool      `json:"state" gorm:"default:'true'"`

	QuizQuestionID uint `json:"quiz_question_id"`
}
