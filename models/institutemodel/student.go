package institutemodel

import (
	"github.com/paulantezana/review/models/reviewmodel"
	"time"
)

// Student struct
type Student struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	DNI      string `json:"dni" gorm:" type:varchar(15); unique; not null"`
	FullName string `json:"full_name" gorm:"type:varchar(250)"`
	Phone    string `json:"phone" gorm:"type:varchar(32)"`
	Gender   string `json:"gender"`

	BirthDate time.Time `json:"birth_date"`

	UserID          uint `json:"user_id"`
	StudentStatusID uint `json:"student_status_id"`

    ProgramID uint `json:"program_id" gorm:"-"` // Only use in JSON

	Reviews []reviewmodel.Review `json:"reviews"`
}
