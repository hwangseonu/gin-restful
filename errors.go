package restful

// HTTPError represents an HTTP error with a status code and message.
// When returned from a handler, the framework automatically sends the
// status code and a JSON body containing the message. The Status field
// is excluded from JSON serialization as it is sent as the HTTP status code.
type HTTPError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Details any    `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	return e.Message
}

// ErrorOption configures optional fields on HTTPError.
type ErrorOption func(*HTTPError)

// WithCode sets an application-specific error code on the HTTPError.
func WithCode(code string) ErrorOption {
	return func(e *HTTPError) {
		e.Code = code
	}
}

// WithDetails sets additional details on the HTTPError.
// Details can be any JSON-serializable value (map, slice, struct, etc.).
func WithDetails(details any) ErrorOption {
	return func(e *HTTPError) {
		e.Details = details
	}
}

// Abort creates an HTTPError with the given status code and message.
// Optional ErrorOption arguments can set Code and Details fields:
//
//	restful.Abort(400, "validation failed",
//	    restful.WithCode("VALIDATION_ERROR"),
//	    restful.WithDetails(map[string]string{"name": "required"}),
//	)
func Abort(status int, message string, opts ...ErrorOption) *HTTPError {
	err := &HTTPError{Status: status, Message: message}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
