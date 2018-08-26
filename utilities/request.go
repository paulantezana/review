package utilities

// Request Type
// 2 = primordial
// 1 = minimal
// 0 = all
type Request struct {
	Search      string `json:"search"`
	CurrentPage uint   `json:"current_page"`
	Limit       uint   `json:"limit"`
	Type        uint   `json:"query"`
}
