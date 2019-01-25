package models

type SubsidiaryUser struct {
	ID           uint `json:"id" gorm:"primary_key"`
	UserID       uint `json:"user_id"`
	SubsidiaryID uint `json:"subsidiary_id"`
	License      bool `json:"license"`
}
