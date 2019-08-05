package models

import "time"

type CourseLevel struct {
    ID          uint      `json:"id" gorm:"primary_key"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    CourseID uint `json:"course_id"`
    Name string `json:"name"`
    Description string `json:"description"`
    ParentId  uint `json:"parent_id"`
}