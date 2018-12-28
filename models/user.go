package models

import (
	"github.com/paulantezana/review/models/institutemodel"
	"github.com/paulantezana/review/models/reviewmodel"
)

// User -- Profiles sa / admin / teacher / secretary
type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	UserName string `json:"user_name" gorm:"type:varchar(64); unique; not null"` //
	Password string `json:"password" gorm:"type:varchar(64); not null"`
	Key      string `json:"key"`
	State    bool   `json:"state" gorm:"default:'true'"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email" gorm:"type:varchar(64)"`
	Freeze   bool   `json:"-"`

	RoleID           uint   `json:"role_id"`
	DefaultProgramID uint   `json:"default_program_id"`
	OldPassword      string `json:"old_password" gorm:"-"`

	Students []institutemodel.Student `json:"students"`
	Teachers []institutemodel.Teacher `json:"teachers"`

	Reviews []reviewmodel.Review `json:"reviews"`
}
