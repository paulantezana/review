package models

import "time"

type CourseNode struct {
    ID          uint      `json:"id" gorm:"primary_key"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    Name string `json:"name"`
    Description string `json:"description"`
    Note float32 `json:"note"`
    Type string `json:"type"`
    
    CourseLevelID uint `json:"course_level_id"`
}