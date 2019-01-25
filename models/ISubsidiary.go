package models

type Subsidiary struct {
	ID                  uint   `json:"id" gorm:"primary_key"`
	Name                string `json:"name"`
	Country             string `json:"country"`
	Department          string `json:"department"`
	Province            string `json:"province"`
	District            string `json:"district"`
	TownCenter          string `json:"town_center"`
	AddressSubsidiary   string `json:"address_subsidiary"`
	RegistrationAddress string `json:"registration_address"`
	Address             string `json:"address"`
	Main                bool   `json:"main"`
	Phone               string `json:"phone"`
}
