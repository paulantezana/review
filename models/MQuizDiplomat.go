package models

import "time"

type QuizDiplomat struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	QuizCount   uint      `json:"quiz_count"`
	CurrentQuiz uint      `json:"current_quiz"`
	Finished    bool      `json:"finished"`
	State       bool      `json:"state"`

	ProgramID uint `json:"program_id"`
}
