package models

type AppModule struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsNew       bool   `json:"is_new"`
}
