package models

import "time"

type StudentProgram struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	StudentID     uint      `json:"student_id"`
	ProgramID     uint      `json:"program_id"`
	ByDefault     bool      `json:"by_default"`
	YearAdmission uint      `json:"year_admission"`
	YearPromotion uint      `json:"year_promotion"`
}
