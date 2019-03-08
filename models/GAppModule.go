package models

import "time"

type AppModules struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ModuleKey string    `json:"module_key"`
	Name      string    `json:"name"`
	Scope     string    `json:"scope"`
	ParentID  uint      `json:"parent_id"`
}
