package models

type AdmissionPayment struct {
	ID          uint    `json:"id" gorm:"primary_key"`
	Payment     float32 `json:"payment"`
	Description string  `json:"description"`

	AdmissionID uint `json:"admission_id"`
}
