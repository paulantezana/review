package librarymodel

type Comment struct {
    ID       uint   `json:"id" gorm:"primary_key"`
    Body string `json:"body"`
    
    UserID uint `json:"user_id"`
    BookID uint `json:"book_id"`
}
