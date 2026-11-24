package response

const (
	// ErrCodeInternalServerError is the error code for internal server errors.
	ErrCodeInternalServerError = "INTERNAL_SERVER_ERROR"

	// ErrCodeServiceUnavailable is the error code for a dependency (e.g. the database) being unreachable.
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

var messages = map[string]string{
	ErrCodeInternalServerError: "Internal server error",
	ErrCodeServiceUnavailable:  "Service unavailable",
}
