package institutemodel

type Semester struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name" gorm:"type:varchar(128); not null"`
	Sequence uint   `json:"sequence"`
	Period   string `json:"period"`
	Year     uint   `json:"year" gorm:"not null"`

	ProgramID uint `json:"program_id"`
}
