package admissionmodel

type Payment struct {
    ID uint `json:"id" gorm:"primary_key"`
    Amount float32 `json:"amount"`
    Reason string `json:"reason"`

    SubsidiaryID uint `json:"subsidiary_id"`
}
