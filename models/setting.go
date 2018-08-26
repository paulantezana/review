package models

type Setting struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	Career       string `json:"company"`
	Manager      string `json:"manager"`
	Email        string `json:"email"`
	Logo         string `json:"logo"`
	CreationDate string `json:"creation_date"`
	ItemTable    uint   `json:"item_table"`
}
