package models

import "time"

type Subsidiary struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Name       string    `json:"name"`
	Country    string    `json:"country"`
	Department string    `json:"department"`
	Province   string    `json:"province"`
	District   string    `json:"district"`
	TownCenter string    `json:"town_center"`
	Address    string    `json:"address"`
	Main       bool      `json:"main"`
	Phone      string    `json:"phone"`

	Programs []Program `json:"programs, omitempty"`
}
