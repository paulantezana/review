package models

import "time"

type MssGroup struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint      `json:"-"`
	UpdateUserId uint      `json:"-"`

	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Avatar        string         `json:"avatar"`
	Date          time.Time      `json:"date"`
	IsActive      bool           `json:"is_active" gorm:"default:'true'"`
	MssUserGroups []MssUserGroup `json:"mss_user_groups"`
}
