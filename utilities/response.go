package utilities

// Response data
type Response struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// ResponsePaginate data pagination
type ResponsePaginate struct {
	Message     string      `json:"message"`
	Success     bool        `json:"success"`
	Data        interface{} `json:"data"`
	Total       uint        `json:"total"`
	CurrentPage uint        `json:"current_page"`
	Limit       uint        `json:"limit"`
}

// Response navigation single
type Navigation struct {
    Title string `json:"title"`
    ID string `json:"id"`
} 
type ResponseNavigation struct {
    Message     string      `json:"message"`
    Success     bool        `json:"success"`
    Data        interface{} `json:"data"`
    Navigation []Navigation `json:"navigation"`
    Total       uint        `json:"total"`
}