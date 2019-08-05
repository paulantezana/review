package models

import "time"

type PreAdmission struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

    School string `json:"school"`
    SchoolPromotionYear uint `json:"school_promotion_year"`

	StudentID          uint `json:"student_id"`
	ProgramID          uint `json:"program_id"`
    AdmissionModalityId uint `json:"admission_modality_id"`
	AdmissionSettingID uint `json:"admission_setting_id"`
}
