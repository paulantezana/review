package models

import "time"

type Course struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name        string    `json:"name"`
	Description string    `json:"description"`
	Poster      string    `json:"poster"`
	Credit      float32   `json:"credit" gorm:"not null"`
	Hours       uint      `json:"hours"  gorm:"not null"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Price       float32   `json:"price"`

	ModuleID   uint `json:"module_id"`
	SemesterID uint `json:"semester_id"`
}
