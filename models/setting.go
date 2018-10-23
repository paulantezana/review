package models

// Setting struct
type Setting struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Logo      string `json:"logo"`
	Institute string `json:"institute"`
	ItemTable uint   `json:"item_table"`
}
