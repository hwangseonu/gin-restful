package gin_restful

import (
	"errors"
	"testing"
)

func TestAbort(t *testing.T) {
	err := Abort(404, "not found")

	if err.Status != 404 {
		t.Errorf("expected status 404, got %d", err.Status)
	}
	if err.Message != "not found" {
		t.Errorf("expected message 'not found', got %q", err.Message)
	}
}

func TestHTTPError_ImplementsError(t *testing.T) {
	var err error = Abort(400, "bad request")
	if err.Error() != "bad request" {
		t.Errorf("expected 'bad request', got %q", err.Error())
	}
}

func TestHTTPError_ErrorsAs(t *testing.T) {
	err := Abort(422, "validation failed")
	var wrapped error = err

	var httpErr *HTTPError
	if !errors.As(wrapped, &httpErr) {
		t.Fatal("errors.As should match *HTTPError")
	}
	if httpErr.Status != 422 {
		t.Errorf("expected status 422, got %d", httpErr.Status)
	}
}
