package models

import "time"

type ModuleSemester struct {
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	SemesterID uint      `json:"semester_id" gorm:"primary_key"`
	ModuleID   uint      `json:"module_id" gorm:"primary_key"`
}
