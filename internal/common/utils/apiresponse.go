package utils

// SuccessResponse ...
type SuccessResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorResponse ...
type ErrorResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func ValidationResponse(error string) ErrorResponse {
	return ErrorResponse{
		Status:  false,
		Message: "Request Failed",
		Error:   error,
	}
}
