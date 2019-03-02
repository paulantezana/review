package models

import "time"

// Program struct
type Program struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name" type:varchar(255); unique; not null"`
	Level     string    `json:"level"`

	SubsidiaryID uint `json:"subsidiary_id"`
}
