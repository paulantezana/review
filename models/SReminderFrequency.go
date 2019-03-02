package models

import "time"

type ReminderFrequency struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Frequency float32   `json:"frequency"`
	IsActive  bool      `json:"is_active"`
}
