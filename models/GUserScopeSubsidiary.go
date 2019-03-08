package models

import "time"

type UserScopeSubsidiary struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
    License   bool      `json:"license"`

    AppModuleID uint `json:"app_module_id"`
	SubsidiaryUserID uint `json:"subsidiary_user_id"`
}
