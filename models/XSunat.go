package models

type Sunat struct {
	RUC              string `json:"ruc"  gorm:"type:varchar(11); unique; not null"`
	SocialReason string `json:"social_reason"`
	State string `json:"state"`
	Location string `json:"location"`
    DomicileCondition string `json:"domicile_condition"`
    FiscalAddress string `json:"fiscal_address"`
}
