package librarymodel

type Like struct {
	ID    uint `json:"id" gorm:"primary_key"`
	Stars uint `json:"stars"`

	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`
}
