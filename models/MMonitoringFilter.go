package models

import "time"

type MonitoringFilter struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Table     string    `json:"table"`
	TableID   uint      `json:"table_id"`
	Type      string    `json:"type"` // subsidiary || program || student || teacher || user || all
}
