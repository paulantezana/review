package messengermodel

import (
    "time"
)

type Session struct {
	ID     uint `json:"id" gorm:"primary_key"`
    IpAddress    string    `json:"ip_address"`
    UserName    string `json:"user_name"`
    LastActivity time.Time `json:"last_activity"`
}
