package models

import "time"

type UserGroup struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Date     time.Time `json:"date"`
    IsActive    bool   `json:"is_active" gorm:"default:'true'"`
	IsAdmin bool `json:"is_admin"`

	UserID  uint `json:"user_id"`
	GroupID uint `json:"group_id"`
}
