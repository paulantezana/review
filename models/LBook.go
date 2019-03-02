package models

import "time"

type BStarts struct {
	UserName string `json:"user_name"`
	Stars    uint8  `json:"stars"`
}

type BookDetail struct {
	HasStart   uint8     `json:"has_start"`   // Current user
	StartValue uint8     `json:"start_value"` // Current user
	Comments   uint      `json:"comments"`    // Comment count
	Views      uint8     `json:"views"`       // Views
	Starts     []BStarts `json:"starts"`      // Other users
}

// Book struct
type Book struct {
	ID               uint       `json:"id" gorm:"primary_key"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Name             string     `json:"name"`
	ShortDescription string     `json:"short_description"`
	LongDescription  string     `json:"long_description"`
	Author           string     `json:"author"`
	Editorial        string     `json:"editorial"`
	YearEdition      uint       `json:"year_edition"`
	Version          uint       `json:"version"`
	EnableDownload   bool       `json:"enable_download"`
	Avatar           string     `json:"avatar"`
	Pdf              string     `json:"pdf"`
	State            bool       `json:"state" gorm:"default:'true'"`
	Views            uint32     `json:"views"`
	Detail           BookDetail `json:"detail" gorm:"-"`

	CategoryID uint `json:"category_id"`
	UserID     uint `json:"user_id"`

	Comments []Comment `json:"comments, omitempty"`
	Readings []Reading `json:"readings, omitempty"`
	Likes    []Like    `json:"likes, omitempty"`
}
