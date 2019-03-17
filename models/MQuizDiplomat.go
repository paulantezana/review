package models

import "time"

type QuizDiplomat struct {
    ID        uint      `json:"id" gorm:"primary_key"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    // ss
}