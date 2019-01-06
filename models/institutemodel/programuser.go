package institutemodel

type ProgramUser struct {
	ID        uint `json:"id" gorm:"primary_key"`
	UserID    uint `json:"user_id"`
	ProgramID uint `json:"program_id"`
	License   bool `json:"license"`
}
