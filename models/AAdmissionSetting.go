package models

import "time"

type AdmissionSetting struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	VacantByProgram uint      `json:"vacant_by_program"`
	Name            string    `json:"name"`
	Year            uint      `json:"year"`
	Seats           uint      `json:"seats"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	PreStartDate    time.Time `json:"pre_start_date"`
	PreEndDate      time.Time `json:"pre_end_date"`
	PreEnabled      bool      `json:"pre_enabled"`
	PreDescription  string    `json:"pre_description"`

	SubsidiaryID uint        `json:"subsidiary_id"`
	Admissions   []Admission `json:"admissions, omitempty"`
}
