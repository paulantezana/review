package models

// Setting struct
type Setting struct {
	ID                         uint   `json:"id" gorm:"primary_key"`
	Logo                       string `json:"logo"`
	Ministry                   string `json:"ministry"`
	Prefix                     string `json:"prefix"`
	PrefixShortName            string `json:"prefix_short_name"`
	Institute                  string `json:"institute"`
	ItemTable                  uint   `json:"item_table"`
	MinHoursPracticePercentage uint   `json:"min_hours_practice_percentage"`
}
