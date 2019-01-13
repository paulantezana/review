package institutemodel

// Program struct
type Program struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name" type:varchar(255); unique; not null"`
	Level string `json:"level"`

	SubsidiaryID uint `json:"subsidiary_id"`
}
