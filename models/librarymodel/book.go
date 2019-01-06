package librarymodel

// Book struct
type Book struct {
	ID               uint   `json:"id" gorm:"primary_key"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	LongDescription  string `json:"long_description"`
	Author           string `json:"author"`
	Editorial        string `json:"editorial"`
	YearEdition      uint   `json:"year_edition"`
	Version          uint   `json:"version"`
	EnableDownload   bool   `json:"enable_download"`
	Avatar           string `json:"avatar"`
	Pdf              string `json:"pdf"`
	State            bool   `json:"state" gorm:"default:'true'"`

	CategoryID uint `json:"category_id"`
	UserID     uint `json:"user_id"`

	Comments []Comment `json:"comments"`
	Readings []Reading `json:"readings"`
	Likes    []Like    `json:"likes"`
}
