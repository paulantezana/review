package models

// Program struct
type Program struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name" type:varchar(255); unique; not null"`

	Students []Student `json:"students, omitempty"`
	Teachers []Teacher `json:"teachers, omitempty"`
	Users    []User    `json:"users, omitempty"`
	Modules  []Module  `json:"modules, omitempty"`
}
