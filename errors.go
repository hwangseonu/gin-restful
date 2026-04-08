package gin_restful

// HTTPError represents an HTTP error with a status code and message.
// When returned from a handler, the framework automatically sends the
// status code and a JSON body containing the message. The Status field
// is excluded from JSON serialization as it is sent as the HTTP status code.
type HTTPError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	return e.Message
}

// Abort creates an HTTPError with the given status code and message.
// Use this in handlers to return structured error responses:
//
//	return nil, 0, restful.Abort(404, "not found")
func Abort(status int, message string) *HTTPError {
	return &HTTPError{Status: status, Message: message}
}
