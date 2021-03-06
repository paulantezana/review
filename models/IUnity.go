package models

import "time"

type Unity struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name" gorm:"type:varchar(128); not null"`
	Credit    float32   `json:"credit" gorm:"not null"`
	Hours     uint      `json:"hours"  gorm:"not null"`
	State     bool      `json:"state" gorm:"default:'true'"`

	ModuleID   uint `json:"module_id"`
	SemesterID uint `json:"semester_id"`
}
