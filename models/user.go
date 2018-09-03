package models

// User -- Profiles sa / admin / teacher / secretary
type User struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	UserName    string `json:"user_name" gorm:"type:varchar(64); unique; not null"`
	Password    string `json:"password" gorm:"type:varchar(64); not null"`
	Profile     string `json:"profile" gorm:"type:varchar(64)"`
	Key         string `json:"key"`
	State       bool   `json:"state" gorm:"default:'true'"`
    Avatar    string `json:"avatar"`
    Email       string `json:"email" gorm:"type:varchar(64)"`

    ProgramID uint `json:"program_id"`
	OldPassword string `json:"old_password" gorm:"-"`

	Reviews []Review `json:"reviews"`
}
