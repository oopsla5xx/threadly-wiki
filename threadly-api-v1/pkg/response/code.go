package response

const (
	// ErrCodeInternalServerError is the error code for internal server errors.
	ErrCodeInternalServerError = "INTERNAL_SERVER_ERROR"
)

var messages = map[string]string{
	ErrCodeInternalServerError: "Internal server error",
}


