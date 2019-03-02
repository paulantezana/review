package models

import "time"

type Semester struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name" gorm:"type:varchar(128); not null"`
	Sequence  uint      `json:"sequence"`
	Period    string    `json:"period"`
	Year      uint      `json:"year" gorm:"not null"`

	ProgramID uint `json:"program_id"`
}
