package models

import "time"

// Category struct
type Category struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ParentID  uint      `json:"parent_id"`
	State     bool      `json:"state" gorm:"default:'true'"`
	ProgramId uint      `json:"program_id"`

	Books []Book `json:"books, omitempty"`
}
