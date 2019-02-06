package models

type AdmissionSetting struct {
	ID              uint   `json:"id" gorm:"primary_key"`
	VacantByProgram uint   `json:"vacant_by_program"`
	Year            uint   `json:"year"`
	Seats           uint   `json:"seats"`
	Description     string `json:"description"`

	Admissions []Admission `json:"admissions, omitempty"`
}
