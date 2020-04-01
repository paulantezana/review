package models

import "time"

type UserRoleModule struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint      `json:"-"`
	UpdateUserId uint      `json:"-"`

	AppModuleID uint `json:"app_module_id"`
	License     bool `json:"license"`
}
