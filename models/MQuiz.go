package models

import "time"

type Quiz struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Message         string    `json:"message"`
	StartDate       time.Time `json:"start_date"`
	StartDateEnable bool      `json:"start_date_enable"`
	EndDate         time.Time `json:"end_date"`
	EndDateEnable   bool      `json:"end_date_enable"`
	LimitTime       uint      `json:"limit_time"`
	LimitTimeFormat uint      `json:"limit_time_format"`
	LimitTimeEnable bool      `json:"limit_time_enable"`
	ShowAnalyze     bool      `json:"show_analyze"`
	BaseNote        uint      `json:"base_note"`
	State           bool      `json:"state" gorm:"default:'true'"`
	Advance         uint      `json:"advance"`
	ProgramID       uint      `json:"program_id"`
	QuizDiplomatID  uint      `json:"quiz_diplomat_id"` // No foreign key

	QuizQuestions []QuizQuestion `json:"quiz_questions"`
}
