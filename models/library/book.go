package library

// Book struct
type Book struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Author         string `json:"author"`
	Editorial      string `json:"editorial"`
	YearEdition    uint   `json:"year_edition"`
	Version        uint   `json:"version"`
	EnableDownload bool   `json:"enable_download"`
	Avatar         string `json:"avatar"`
	Pdf            string `json:"pdf"`
	State          bool   `json:"state" gorm:"default:'true'"`

	CategoryID uint `json:"category_id"`
}
