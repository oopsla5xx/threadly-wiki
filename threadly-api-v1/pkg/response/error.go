package response

import "net/http"

// AppError represents an application error with an HTTP status code,
// a machine-readable error code, and a client-facing message.
type AppError struct {
	StatusCode int
	Code       string
	Message    string
}

func (applicationError *AppError) Error() string {
	return applicationError.Message
}

// NewAppError creates an application error with the provided status,
// error code, and message.
func NewAppError(
	statusCode int,
	code string,
	message string,
) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// InternalServerError creates a standardized HTTP 500 application error.
func InternalServerError() *AppError {
	return NewAppError(
		http.StatusInternalServerError,
		ErrCodeInternalServerError,
		messages[ErrCodeInternalServerError],
	)
}

// ServiceUnavailable creates a standardized HTTP 503 application error.
func ServiceUnavailable() *AppError {
	return NewAppError(
		http.StatusServiceUnavailable,
		ErrCodeServiceUnavailable,
		messages[ErrCodeServiceUnavailable],
	)
}
