package models

import "time"

type MonitoringFilterDetail struct {
	ID                 uint      `json:"id" gorm:"primary_key"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	MonitoringFilterID uint      `json:"monitoring_filter_id"`
	ReferenceID        uint      `json:"reference_id"`
	ReferenceName      string    `json:"reference_name"`
}
