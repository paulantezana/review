package models

// Category struct
type Category struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	ParentID uint   `json:"parent_id"`
	State    bool   `json:"state" gorm:"default:'true'"`

	Books    []Book     `json:"books"`
	Children []Category `json:"children"`
}
