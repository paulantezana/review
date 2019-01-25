package models

// Module struct
type Module struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Sequence    uint   `json:"sequence"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Points      uint   `json:"points"`
	Hours       uint   `json:"hours"`

	ProgramID uint `json:"program_id"`

	Semesters []ModuleSemester `json:"semesters"`
}
