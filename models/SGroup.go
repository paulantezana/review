package models

import "time"

type Group struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Name     string    `json:"name"`
	Avatar   string    `json:"avatar"`
	Date     time.Time `json:"date"`
	IsActive bool      `json:"is_active" gorm:"default:'true'"`

	UserGroups []UserGroup `json:"user_groups"`
}
