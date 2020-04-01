package models

import "time"

type MssReminderFrequency struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint      `json:"-"`
	UpdateUserId uint      `json:"-"`

	Name      string  `json:"name"`
	Frequency float32 `json:"frequency"`
	IsActive  bool    `json:"is_active"`
}
