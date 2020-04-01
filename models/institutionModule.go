package models

import "time"

type InstitutionModule struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AppModuleID uint      `json:"app_module_id"`
	State       bool      `json:"state"`
}
