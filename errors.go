package gin_restful

type HTTPError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func Abort(status int, message string) *HTTPError {
	return &HTTPError{Status: status, Message: message}
}
