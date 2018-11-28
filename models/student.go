package models

import "time"

// Student struct
type Student struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	DNI      string `json:"dni" gorm:" type:varchar(15); unique; not null"`
	FullName string `json:"full_name" gorm:"type:varchar(250)"`
	Phone    string `json:"phone" gorm:"type:varchar(32)"`
	State    bool   `json:"state" gorm:"default:'true'"`
	Gender   string `json:"gender"`

	BirthDate     time.Time `json:"birth_date"`
	AdmissionYear uint      `json:"admission_year"`
	PromotionYear uint      `json:"promotion_year"`

	ProgramID uint `json:"program_id"`
	UserID    uint `json:"user_id"`

	Reviews []Review `json:"reviews"`
	User    User     `json:"user"`
}
