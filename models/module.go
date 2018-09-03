package models

type Module struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Sequence    uint   `json:"sequence"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Points      uint   `json:"points"`
	Hours       uint   `json:"hours"`
	Semester    string `json:"semester"`

	ProgramID uint `json:"program_id"`
	
	Reviews []Review `json:"reviews"`
}
