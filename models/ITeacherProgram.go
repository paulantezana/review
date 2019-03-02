package models

import "time"

type TeacherProgram struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TeacherID uint      `json:"teacher_id"`
	ProgramID uint      `json:"program_id"`
	ByDefault bool      `json:"by_default"`
	Type      string    `json:"type"` // cross // career
}
