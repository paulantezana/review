package models

// User -- Profiles sa / admin / teacher / secretary
type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	UserName string `json:"user_name" gorm:"type:varchar(64); unique; not null"` //
	Password string `json:"password, omitempty" gorm:"type:varchar(64); not null"`
	Key      string `json:"key, omitempty"`
	State    bool   `json:"state" gorm:"default:'true'"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email" gorm:"type:varchar(64)"`
	Freeze   bool   `json:"-"`

	RoleID      uint   `json:"role_id"`
	OldPassword string `json:"old_password" gorm:"-"`

	Students []Student `json:"students,omitempty"`
	Teachers []Teacher `json:"teachers,omitempty"`
	Reviews  []Review  `json:"reviews,omitempty"`
	Comments []Comment `json:"coos,omitempty"`
}
