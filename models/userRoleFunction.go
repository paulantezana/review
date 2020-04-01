package models

import "time"

type UserRoleFunction struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint      `json:"-"`
	UpdateUserId uint      `json:"-"`

	UserRoleModuleID    uint `json:"user_role_module_id"`
	AppModuleFunctionID uint `json:"app_module_function_id"`
	License             bool `json:"license"`
}
