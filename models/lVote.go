package models

type Vote struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CommentID uint `json:"comment_id" gorm:"not null"`
	UserID    uint `json:"user_id" gorm:"not null"`
	Value     bool `json:"value" gorm:"not null"`
}
