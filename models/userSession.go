package models

import (
	"time"
)

type UserSession struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IpAddress string    `json:"ip_address"`

	UserID       uint      `json:"user_id"`
	IsOnline     uint      `json:"is_online"`
	LastActivity time.Time `json:"last_activity"`
}
