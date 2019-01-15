package reviewmodel

import (
	"time"
)

// Review struct
type Review struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	ApprobationDate time.Time `json:"approbation_date"`

	ModuleId         uint `json:"module_id"`
	StudentProgramID uint `json:"student_program_id"`
	CreatorID        uint `json:"creator_id"`
	TeacherID        uint `json:"teacher_id"`

	ReviewDetails []ReviewDetail `json:"review_details"`
}
