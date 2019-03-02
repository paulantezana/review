package models

import "time"

type AdmissionPayment struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Payment     float32   `json:"payment"`
	Description string    `json:"description"`

	AdmissionID uint `json:"admission_id"`
}
