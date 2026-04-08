package gin_restful

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type testBody struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age"`
}

func TestBind_Success(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"alice","age":30}`))
	c.Request.Header.Set("Content-Type", "application/json")

	body, err := Bind[testBody](c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body.Name != "alice" {
		t.Errorf("expected name 'alice', got %q", body.Name)
	}
	if body.Age != 30 {
		t.Errorf("expected age 30, got %d", body.Age)
	}
}

func TestBind_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{invalid`))
	c.Request.Header.Set("Content-Type", "application/json")

	_, err := Bind[testBody](c)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestBind_ValidationError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"age":25}`))
	c.Request.Header.Set("Content-Type", "application/json")

	_, err := Bind[testBody](c)
	if err == nil {
		t.Fatal("expected validation error for missing required field 'name'")
	}
}

func TestMustBind_Success(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"bob","age":25}`))
	c.Request.Header.Set("Content-Type", "application/json")

	body := MustBind[testBody](c)
	if body == nil {
		t.Fatal("expected non-nil body")
	}
	if body.Name != "bob" {
		t.Errorf("expected name 'bob', got %q", body.Name)
	}
	if c.IsAborted() {
		t.Error("should not be aborted on success")
	}
}

func TestMustBind_Failure_AbortsWith400(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"age":25}`))
	c.Request.Header.Set("Content-Type", "application/json")

	body := MustBind[testBody](c)
	if body != nil {
		t.Fatal("expected nil body on failure")
	}
	if !c.IsAborted() {
		t.Error("expected context to be aborted")
	}
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
