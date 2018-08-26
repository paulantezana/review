package models

type Student struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	DNI      string `json:"dni" gorm:" type:varchar(15); unique; not null"`
	FullName string `json:"full_name" gorm:"type:varchar(250)"`
	Email    string `json:"email" gorm:"type:varchar(64)"`
	Phone    string `json:"phone" gorm:"type:varchar(32)"`
	State    bool   `json:"state" gorm:"default:'true'"`

	Reviews []Review `json:"reviews"`
}
