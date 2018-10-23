package models

import "time"

// Review struct
type Review struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	Module          string    `json:"module"`
	Semester        string    `json:"semester"`
	Supervisor      string    `json:"supervisor"`
	ApprobationDate time.Time `json:"approbation_date"`

	ModuleId  uint `json:"module_id"`
	StudentID uint `json:"student_id"`
	UserID    uint `json:"user_id"`

	ReviewDetails []ReviewDetail `json:"review_details"`
}
