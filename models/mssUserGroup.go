package models

import "time"

type MssUserGroup struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint      `json:"-"`
	UpdateUserId uint      `json:"-"`

	IsActive   bool `json:"is_active" gorm:"default:'true'"`
	IsAdmin    bool `json:"is_admin"`
	UserID     uint `json:"user_id"`
	MssGroupID uint `json:"mss_group_id"`
}
