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

// ResponseScroll data scroll
type ResponseScroll struct {
    Message     string      `json:"message"`
    Success     bool        `json:"success"`
    Data        interface{} `json:"data"`
    HasMore bool `json:"has_more"`
}

// SocketResponse
type SocketResponse struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
