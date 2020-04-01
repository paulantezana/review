package models

import "time"

type App struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	LastUpdate time.Time `json:"last_update"`
}
