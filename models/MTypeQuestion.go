package models

// TypeQuestion struct
type TypeQuestion struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`

	Questions []Question `json:"questions, omitempty"`
}
