package monitoringmodel

import (
	"time"
)

// Poll struct
type Poll struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Message     string    `json:"message"`
	Weather     bool      `json:"weather"` //definite / undefined
	State       bool      `json:"state" gorm:"default:'true'"`

	ProgramID uint `json:"program_id"`

	Questions []Question `json:"questions, omitempty"`
}
