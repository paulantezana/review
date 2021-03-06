package models

import "time"

// ReviewDetail struct
type ReviewDetail struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Hours     uint      `json:"hours"`
	Note      uint      `json:"note"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	ReviewID  uint `json:"review_id"`
	CompanyID uint `json:"company_id"`
}
