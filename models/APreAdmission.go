package models

import "time"

type PreAdmission struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	StudentID          uint `json:"student_id"`
	ProgramID          uint `json:"program_id"`
	AdmissionSettingID uint `json:"admission_setting_id"`
}
