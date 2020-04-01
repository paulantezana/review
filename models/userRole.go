package models

import "time"

type UserRole struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint      `json:"-"`
	UpdateUserId uint      `json:"-"`

	IsMain bool   `json:"is_main"`
	ParentID uint `json:"parent_id"`
	Name   string `json:"name"`
	State  bool   `json:"state"`
}
