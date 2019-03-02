package models

import (
	"time"
)

type Admission struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Observation   string    `json:"observation"`
	Exonerated    bool      `json:"exonerated"`
	ExamNote      float32   `json:"exam_note"`
	ExamDate      time.Time `json:"exam_date"`
	AdmissionDate time.Time `json:"admission_date"`
	Year          uint      `json:"year"`
	Classroom     uint      `json:"classroom"`
	Seat          uint      `json:"seat"`
	State         bool      `json:"state" gorm:"default:'true'"`

	StudentID          uint `json:"student_id"`
	ProgramID          uint `json:"program_id"`
	UserID             uint `json:"user_id"`
	AdmissionSettingID uint `json:"admission_setting_id"`
}
