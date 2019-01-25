package models

type Payment struct {
	ID       uint    `json:"id" gorm:"primary_key"`
	Amount   float32 `json:"amount"`
	Reason   string  `json:"reason"`
	Category string  `json:"category"` // admission || enrollment
}
