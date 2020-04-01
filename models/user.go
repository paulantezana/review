package models

import "time"

// User -- Profiles sa / admin / teacher / secretary
type User struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreateUserId uint      `json:"-"`
	UpdateUserId uint      `json:"-"`

	UserName     string `json:"user_name" gorm:"type:varchar(64); unique; not null"` //
	Password     string `json:"password, omitempty" gorm:"type:varchar(64); not null"`
	TempKey      string `json:"temp_key, omitempty"`
	QrKey        string `json:"qr_key"`
	BarCodeKey   string `json:"bar_code_key"`
	BiometricKey string `json:"biometric_key"`
	State        bool   `json:"state" gorm:"default:'true'"`
	Avatar       string `json:"avatar"`
	Email        string `json:"email" gorm:"type:varchar(64)"`
	Freeze       bool   `json:"-"`

	UserRoleID  uint   `json:"user_role_id"`
	InstitutionID uint `json:"institution_id"`
	OldPassword string `json:"old_password" gorm:"-"`
}
