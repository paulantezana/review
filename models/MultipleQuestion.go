package models

// MultipleQuestion struct
type MultipleQuestion struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Label string `json:"label"`
	State bool   `json:"state" gorm:"default:'true'"`

	QuestionID uint `json:"question_id"`
}
