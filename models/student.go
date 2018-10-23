package models

import "time"

// Student struct
type Student struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	DNI      string `json:"dni" gorm:" type:varchar(15); unique; not null"`
	FullName string `json:"full_name" gorm:"type:varchar(250)"`
	Phone    string `json:"phone" gorm:"type:varchar(32)"`
	State    bool   `json:"state" gorm:"default:'true'"`

	BirthDate     time.Time `json:"birth_date"`
	AdmissionDate time.Time `json:"admission_date"`
	PromotionDate time.Time `json:"promotion_date"`

	ProgramID uint `json:"program_id"`

	Reviews []Review `json:"reviews"`
	User    User     `json:"user"`
}
