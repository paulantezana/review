package models

import "time"

// Company struct
type Company struct {
	ID               uint           `json:"id" gorm:"primary_key"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	RUC              string         `json:"ruc"  gorm:"type:varchar(11); unique; not null"`
	NameSocialReason string         `json:"name_social_reason"`
	Address          string         `json:"address"`
	Manager          string         `json:"manager"`
	Phone            string         `json:"phone"`
	CompanyType      string         `json:"company_type"` // 1 = public || 2 = private
	ReviewDetails    []ReviewDetail `json:"review_details"`
}

// TableName function table rename
func (Company) TableName() string {
	return "companies"
}
