package models

import "time"

// TypeQuestion struct
type TypeQuestion struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`

	//Questions []Question `json:"questions, omitempty"`
}
