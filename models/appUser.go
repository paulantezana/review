package models

import "time"

type AppUser struct {
    ID           uint      `json:"id" gorm:"primary_key"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`

    UserName     string `json:"user_name" gorm:"type:varchar(64); unique; not null"` //
    Password     string `json:"password, omitempty" gorm:"type:varchar(64); not null"`
    TempKey      string `json:"temp_key, omitempty"`
    State        bool   `json:"state" gorm:"default:'true'"`
    Avatar       string `json:"avatar"`
    Email        string `json:"email" gorm:"type:varchar(64)"`
    OldPassword string `json:"old_password" gorm:"-"`
}
