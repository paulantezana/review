package monitoring

type TypeQuestion struct {
    ID    uint   `json:"id" gorm:"primary_key"`
    Name  string `json:"name"`
    State bool   `json:"state" gorm:"default:'true'"`

    Questions []Question `json:"questions"`
}
