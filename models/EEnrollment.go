package models

import "time"

type Enrollment struct {
    ID        uint      `json:"id" gorm:"primary_key"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`

    StudentID uint `json:"student_id"`
    Observation string `json:"observation"`
}
