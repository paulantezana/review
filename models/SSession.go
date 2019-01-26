package models

import (
	"time"
)

type Session struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	IpAddress    string    `json:"ip_address"`
	UserID       uint      `json:"user_id"`
	LastActivity time.Time `json:"last_activity"`
}
