package restful

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

func TestAbort_WithCode(t *testing.T) {
	err := Abort(400, "bad request", WithCode("VALIDATION_ERROR"))
	if err.Code != "VALIDATION_ERROR" {
		t.Errorf("expected code 'VALIDATION_ERROR', got %q", err.Code)
	}
}

func TestAbort_WithDetails(t *testing.T) {
	details := map[string]string{"name": "required"}
	err := Abort(400, "validation failed", WithDetails(details))
	if err.Details == nil {
		t.Fatal("expected non-nil details")
	}
	d, ok := err.Details.(map[string]string)
	if !ok {
		t.Fatal("expected details to be map[string]string")
	}
	if d["name"] != "required" {
		t.Errorf("expected details[name]='required', got %q", d["name"])
	}
}

func TestAbort_WithCodeAndDetails(t *testing.T) {
	err := Abort(422, "invalid",
		WithCode("INVALID_INPUT"),
		WithDetails([]string{"field1", "field2"}),
	)
	if err.Code != "INVALID_INPUT" {
		t.Errorf("expected code 'INVALID_INPUT', got %q", err.Code)
	}
	if err.Details == nil {
		t.Fatal("expected non-nil details")
	}
}

func TestAbort_BackwardsCompatible(t *testing.T) {
	err := Abort(404, "not found")
	if err.Code != "" {
		t.Errorf("expected empty code, got %q", err.Code)
	}
	if err.Details != nil {
		t.Errorf("expected nil details, got %v", err.Details)
	}
}
