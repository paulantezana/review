package models

import "time"

type Admission struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	Observation   string    `json:"observation"`
	Exonerated    bool      `json:"exonerated"`
	ExamNote      float32   `json:"exam_note"`
	ExamDate      time.Time `json:"exam_date"`
	AdmissionDate time.Time `json:"admission_date"`
	Year          uint      `json:"year"`
	Classroom     uint      `json:"classroom"`
	Seat          uint      `json:"seat"`

	StudentID uint `json:"student_id"`
	ProgramID uint `json:"program_id"`
	UserID    uint `json:"user_id"`

	State bool `json:"state" gorm:"default:'true'"`
}
