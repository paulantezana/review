package models

import "time"

type AdmissionSetting struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	VacantByProgram uint      `json:"vacant_by_program"`
	Year            uint      `json:"year"`
	Seats           uint      `json:"seats"`
	Description     string    `json:"description"`
	ShowInWeb       bool      `json:"show_in_web"`

	SubsidiaryID uint        `json:"subsidiary_id"`
	Admissions   []Admission `json:"admissions, omitempty"`
}
