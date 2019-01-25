package models

type ModuleSemester struct {
	SemesterID uint `json:"semester_id" gorm:"primary_key"`
	ModuleID   uint `json:"module_id" gorm:"primary_key"`
}
