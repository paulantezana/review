package models

import "time"

// Question struct
type Question struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Required  bool      `json:"required"`
	Position  uint      `json:"position"`
	State     bool      `json:"state" gorm:"default:'true'"`

	TypeQuestionID uint `json:"type_question_id"`
	PollID         uint `json:"poll_id"`

	MultipleQuestions []MultipleQuestion `json:"multiple_questions"`
	AnswerDetails     []AnswerDetail     `json:"answer_details"`
}
