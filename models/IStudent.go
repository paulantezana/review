package models

import (
	"time"
)

// Student struct
type Student struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	DNI         string    `json:"dni" gorm:" type:varchar(15); unique; not null"`
	FullName    string    `json:"full_name" gorm:"type:varchar(250)"`
	Phone       string    `json:"phone" gorm:"type:varchar(32)"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	BirthDate   time.Time `json:"birth_date"`
	BirthPlace  string    `json:"birth_place"`
	Country     string    `json:"country"`
	District    string    `json:"district"`
	Province    string    `json:"province"`
	Region      string    `json:"region"`
	MarketStall string    `json:"market_stall"`
	CivilStatus string    `json:"civil_status"`
	IsWork      string    `json:"is_work"` // y yes || n = no

	UserID          uint `json:"user_id"`
	StudentStatusID uint `json:"student_status_id"`

	ProgramID uint `json:"program_id" gorm:"-"` // Only use in JSON

	Reviews []Review `json:"reviews"`
}
