package models

import "time"

type EnrollmentPayment struct {
    ID        uint      `json:"id" gorm:"primary_key"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    EnrollmentID uint `json:"enrollment_id"`
    Observation string `json:"observation"`
    Amount float32 `json:"amount"`
    Reason string `json:"reason"`
}