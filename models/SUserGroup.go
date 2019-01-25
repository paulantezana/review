package models

import "time"

type UserGroup struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Date     time.Time `json:"date"`
	IsActive bool      `json:"is_active"`

	UserID  uint `json:"user_id"`
	GroupID uint `json:"group_id"`
}
