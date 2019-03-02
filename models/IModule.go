package models

import "time"

// Module struct
type Module struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Sequence    uint      `json:"sequence"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Points      uint      `json:"points"`
	Hours       uint      `json:"hours"`

	ProgramID uint `json:"program_id"`

	Semesters []ModuleSemester `json:"semesters"`
}
