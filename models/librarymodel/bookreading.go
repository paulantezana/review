package librarymodel

// Reading struct
type Reading struct {
	ID   uint `json:"id" gorm:"primary_key"`
	Date uint `json:"date"`

	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`
}
