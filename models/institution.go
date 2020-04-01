package models

import "time"

// Setting struct
type Institution struct {
	ID                         uint      `json:"id" gorm:"primary_key"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
	Logo                       string    `json:"logo"`
	Ministry                   string    `json:"ministry"`
	NationalEmblem             string    `json:"national_emblem"`
	Prefix                     string    `json:"prefix"`
	PrefixShortName            string    `json:"prefix_short_name"`
	Institute                  string    `json:"institute"`
	MinHoursPracticePercentage uint      `json:"min_hours_practice_percentage"`
	YearName                   string    `json:"year_name"`

	Director                   string `json:"director"`
	AcademicLevelDirector      string `json:"academic_level_director"`
	ShortAcademicLevelDirector string `json:"short_academic_level_director"`
	ResolutionAuthorization    string `json:"resolution_authorization"`
	ResolutionRenovation       string `json:"resolution_renovation"`
	ModularCode                string `json:"modular_code"`
	ManagementType             string `json:"management_type"`
	DreGre                     string `json:"dre_gre"`

	WebSite string `json:"web_site"`
	Email   string `json:"email"`
}
