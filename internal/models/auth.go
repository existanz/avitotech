package models

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// NewErrorResponse creates a new ErrorResponse instance.
func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Error: err.Error()}
}
