package models

// Setting struct
type Setting struct {
	ID                         uint   `json:"id" gorm:"primary_key"`
	Logo                       string `json:"logo"`
	Ministry                   string `json:"ministry"`
	Prefix                     string `json:"prefix"`
	PrefixShortName            string `json:"prefix_short_name"`
	Institute                  string `json:"institute"`
	Director                   string `json:"director"`
	AcademicLevelDirector      string `json:"academic_level_director"`
	ShortAcademicLevelDirector string `json:"short_academic_level_director"`
	MinHoursPracticePercentage uint   `json:"min_hours_practice_percentage"`
}
