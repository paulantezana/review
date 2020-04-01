package models

import "time"

type BStarts struct {
	UserName string `json:"user_name"`
	Stars    uint8  `json:"stars"`
}

type PostDetail struct {
	HasStart   uint8     `json:"has_start"`   // Current user
	StartValue uint8     `json:"start_value"` // Current user
	Comments   uint      `json:"comments"`    // Comment count
	Views      uint8     `json:"views"`       // Views
	Starts     []BStarts `json:"starts"`      // Other users
}

// Book struct
type Post struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint `json:"-"`
	UpdateUserId uint `json:"-"`

	Title            string     `json:"title"`
	ShortDescription string     `json:"short_description"`
	LongDescription  string     `json:"long_description"`
	Avatar           string     `json:"avatar"`
	State            bool       `json:"state" gorm:"default:'true'"`
	Views            uint32     `json:"views"`
	Detail           PostDetail `json:"detail" gorm:"-"`

	UserID uint `json:"user_id"`
	PostCategoryID uint `json:"post_category_id"`
	PostTypeID uint `json:"post_type_id"`
	Comments []PostComment `json:"comments, omitempty"`
	Readings []PostReading `json:"readings, omitempty"`
	Likes    []PostLike    `json:"likes, omitempty"`
}
