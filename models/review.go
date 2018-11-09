package models

import "time"

// Review struct
type Review struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	ApprobationDate time.Time `json:"approbation_date"`

	ModuleId  uint `json:"module_id"`
	StudentID uint `json:"student_id"`
	UserID    uint `json:"user_id"`
	TeacherID uint `json:"teacher_id"`

	ReviewDetails []ReviewDetail `json:"review_details"`
}
