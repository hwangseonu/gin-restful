package gin_restful

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// --- test resources ---

type readOnlyResource struct{}

func (r *readOnlyResource) List(c *gin.Context) (any, int, error) {
	return []string{"a", "b"}, http.StatusOK, nil
}

func (r *readOnlyResource) Get(id string, c *gin.Context) (any, int, error) {
	if id == "1" {
		return gin.H{"id": "1", "name": "item1"}, http.StatusOK, nil
	}
	return nil, 0, Abort(404, "not found")
}

type createOnlyResource struct{}

func (r *createOnlyResource) Post(c *gin.Context) (any, int, error) {
	body, err := Bind[testBody](c)
	if err != nil {
		return nil, 0, Abort(400, err.Error())
	}
	return gin.H{"name": body.Name}, http.StatusCreated, nil
}

type fullCRUDResource struct{}

func (r *fullCRUDResource) List(c *gin.Context) (any, int, error) {
	return []string{}, http.StatusOK, nil
}

func (r *fullCRUDResource) Get(id string, c *gin.Context) (any, int, error) {
	return gin.H{"id": id}, http.StatusOK, nil
}

func (r *fullCRUDResource) Post(c *gin.Context) (any, int, error) {
	return gin.H{"created": true}, http.StatusCreated, nil
}

func (r *fullCRUDResource) Put(id string, c *gin.Context) (any, int, error) {
	return gin.H{"updated": id}, http.StatusOK, nil
}

func (r *fullCRUDResource) Patch(id string, c *gin.Context) (any, int, error) {
	return gin.H{"patched": id}, http.StatusOK, nil
}

func (r *fullCRUDResource) Delete(id string, c *gin.Context) (any, int, error) {
	return nil, http.StatusNoContent, nil
}

type errorResource struct{}

func (r *errorResource) Get(id string, c *gin.Context) (any, int, error) {
	return nil, 0, Abort(404, "not found")
}

// --- helpers ---

func setupRouter(path string, resource any) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := NewAPI(engine, "/api")
	api.AddResource(path, resource)
	return engine
}

func doRequest(engine *gin.Engine, method, path string, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

// --- tests: partial implementation ---

func TestAddResource_ReadOnly_RegistersOnlyGetRoutes(t *testing.T) {
	engine := setupRouter("/items", &readOnlyResource{})

	// GET /api/items should work
	w := doRequest(engine, "GET", "/api/items", "")
	if w.Code != http.StatusOK {
		t.Errorf("GET /api/items: expected 200, got %d", w.Code)
	}

	// GET /api/items/1 should work
	w = doRequest(engine, "GET", "/api/items/1", "")
	if w.Code != http.StatusOK {
		t.Errorf("GET /api/items/1: expected 200, got %d", w.Code)
	}

	// POST /api/items should NOT be registered (404 or 405)
	w = doRequest(engine, "POST", "/api/items", `{"name":"test"}`)
	if w.Code != http.StatusNotFound && w.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /api/items: expected 404/405, got %d", w.Code)
	}

	// PUT /api/items/1 should NOT be registered
	w = doRequest(engine, "PUT", "/api/items/1", `{"name":"test"}`)
	if w.Code != http.StatusNotFound && w.Code != http.StatusMethodNotAllowed {
		t.Errorf("PUT /api/items/1: expected 404/405, got %d", w.Code)
	}

	// DELETE /api/items/1 should NOT be registered
	w = doRequest(engine, "DELETE", "/api/items/1", "")
	if w.Code != http.StatusNotFound && w.Code != http.StatusMethodNotAllowed {
		t.Errorf("DELETE /api/items/1: expected 404/405, got %d", w.Code)
	}
}

func TestAddResource_CreateOnly_RegistersOnlyPostRoute(t *testing.T) {
	engine := setupRouter("/items", &createOnlyResource{})

	// POST should work
	w := doRequest(engine, "POST", "/api/items", `{"name":"alice","age":30}`)
	if w.Code != http.StatusCreated {
		t.Errorf("POST /api/items: expected 201, got %d", w.Code)
	}

	// GET should NOT be registered
	w = doRequest(engine, "GET", "/api/items", "")
	if w.Code != http.StatusNotFound && w.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET /api/items: expected 404/405, got %d", w.Code)
	}
}

// --- tests: full CRUD ---

func TestAddResource_FullCRUD(t *testing.T) {
	engine := setupRouter("/items", &fullCRUDResource{})

	tests := []struct {
		method string
		path   string
		body   string
		status int
	}{
		{"GET", "/api/items", "", 200},
		{"GET", "/api/items/1", "", 200},
		{"POST", "/api/items", `{}`, 201},
		{"PUT", "/api/items/1", `{}`, 200},
		{"PATCH", "/api/items/1", `{}`, 200},
		{"DELETE", "/api/items/1", "", 204},
	}

	for _, tt := range tests {
		w := doRequest(engine, tt.method, tt.path, tt.body)
		if w.Code != tt.status {
			t.Errorf("%s %s: expected %d, got %d", tt.method, tt.path, tt.status, w.Code)
		}
	}
}

// --- tests: error handling ---

func TestAddResource_HTTPError_ReturnsCorrectStatusAndBody(t *testing.T) {
	engine := setupRouter("/items", &errorResource{})

	w := doRequest(engine, "GET", "/api/items/999", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp["message"] != "not found" {
		t.Errorf("expected message 'not found', got %q", resp["message"])
	}
}

// --- tests: Bind integration ---

func TestAddResource_BindValidationError_Returns400(t *testing.T) {
	engine := setupRouter("/items", &createOnlyResource{})

	// missing required field "name"
	w := doRequest(engine, "POST", "/api/items", `{"age":25}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- tests: 204 No Content has no body ---

func TestAddResource_NoContent_EmptyBody(t *testing.T) {
	engine := setupRouter("/items", &fullCRUDResource{})

	w := doRequest(engine, "DELETE", "/api/items/1", "")
	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Errorf("expected empty body for 204, got %q", w.Body.String())
	}
}

// --- tests: invalid resource panics ---

func TestAddResource_NoInterface_Panics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for resource implementing no interfaces")
		}
	}()

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := NewAPI(engine, "/api")
	api.AddResource("/items", &struct{}{}) // implements nothing
}

// --- tests: non-HTTPError hides internal message ---

type internalErrorResource struct{}

func (r *internalErrorResource) Get(id string, c *gin.Context) (any, int, error) {
	return nil, 0, fmt.Errorf("db connection failed: host=secret-db.internal")
}

func TestAddResource_InternalError_HidesMessage(t *testing.T) {
	engine := setupRouter("/items", &internalErrorResource{})

	w := doRequest(engine, "GET", "/api/items/1", "")
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp["message"] != "internal server error" {
		t.Errorf("expected generic message, got %q (internal details leaked)", resp["message"])
	}
}

// --- tests: error handler aborts middleware chain ---

func TestAddResource_Error_AbortsMiddlewareChain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	middlewareRan := false
	engine.Use(func(c *gin.Context) {
		c.Next()
		// this runs after the handler; check if aborted
		if !c.IsAborted() {
			middlewareRan = true
		}
	})

	api := NewAPI(engine, "/api")
	api.AddResource("/items", &errorResource{})

	doRequest(engine, "GET", "/api/items/1", "")
	if middlewareRan {
		t.Error("middleware should not have continued after error (c.IsAborted should be true)")
	}
}

// --- tests: path normalization ---

func TestNewAPI_PathNormalization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := NewAPI(engine, "/api/")
	api.AddResource("/items", &readOnlyResource{})

	// Should work without double slash issue
	w := doRequest(engine, "GET", "/api/items", "")
	if w.Code != http.StatusOK {
		t.Errorf("GET /api/items: expected 200, got %d (path normalization issue)", w.Code)
	}
}

// --- tests: custom error handler ---

func TestWithErrorHandler_CustomFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := NewAPI(engine, "/api",
		WithErrorHandler(func(c *gin.Context, err error, status int) {
			var httpErr *HTTPError
			if errors.As(err, &httpErr) {
				c.AbortWithStatusJSON(httpErr.Status, gin.H{
					"error": gin.H{
						"message": httpErr.Message,
						"code":    httpErr.Code,
					},
				})
				return
			}
			c.AbortWithStatusJSON(500, gin.H{"error": gin.H{"message": "server error"}})
		}),
	)
	api.AddResource("/items", &errorResource{})

	w := doRequest(engine, "GET", "/api/items/1", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	var resp map[string]map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp["error"]["message"] != "not found" {
		t.Errorf("expected custom error format with 'not found', got %v", resp)
	}
}

func TestWithErrorHandler_DefaultWhenNotSet(t *testing.T) {
	engine := setupRouter("/items", &errorResource{})

	w := doRequest(engine, "GET", "/api/items/1", "")
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp["message"] != "not found" {
		t.Errorf("expected default format with 'message', got %v", resp)
	}
}

// --- tests: HTTPError with Code and Details in response ---

type detailedErrorResource struct{}

func (r *detailedErrorResource) Get(id string, c *gin.Context) (any, int, error) {
	return nil, 0, Abort(422, "validation failed",
		WithCode("VALIDATION_ERROR"),
		WithDetails(gin.H{"field": "name is required"}),
	)
}

func TestHTTPError_CodeAndDetails_InResponse(t *testing.T) {
	engine := setupRouter("/items", &detailedErrorResource{})

	w := doRequest(engine, "GET", "/api/items/1", "")
	if w.Code != 422 {
		t.Errorf("expected 422, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp["code"] != "VALIDATION_ERROR" {
		t.Errorf("expected code 'VALIDATION_ERROR', got %v", resp["code"])
	}
	if resp["details"] == nil {
		t.Error("expected non-nil details in response")
	}
}

func TestNewAPI_TrailingSlashNormalization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := NewAPI(engine, "/api")
	api.AddResource("/items/", &readOnlyResource{})

	// Trailing slash should be stripped — no 301 redirect
	w := doRequest(engine, "GET", "/api/items", "")
	if w.Code != http.StatusOK {
		t.Errorf("GET /api/items: expected 200, got %d (trailing slash issue)", w.Code)
	}
}
