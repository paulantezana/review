package models

type AppModuleFunction struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	AppModuleID uint   `json:"app_module_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsNew       bool   `json:"is_new"`
}
