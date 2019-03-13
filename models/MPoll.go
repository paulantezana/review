package models

import (
	"time"
)

// Poll struct
type Poll struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Message     string    `json:"message"`
    StartDate   time.Time `json:"start_date"`
    StartDateEnable bool `json:"start_date_enable"`
    EndDate     time.Time `json:"end_date"`
    EndDateEnable bool `json:"end_date_enable"`
	State       bool      `json:"state" gorm:"default:'true'"`
    ShowAnalyze bool `json:"show_analyze"`
	ProgramID uint `json:"program_id"`

	Questions []Question `json:"questions, omitempty"`
}
