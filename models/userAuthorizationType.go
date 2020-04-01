package models

type UserAuthorizationType struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	State bool   `json:"state"`
}
