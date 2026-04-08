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

// --- BindQuery tests ---

type testQuery struct {
	Page    int    `form:"page" binding:"gte=1"`
	PerPage int    `form:"per_page"`
	Search  string `form:"search"`
}

func TestBindQuery_Success(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=2&per_page=20&search=hello", nil)

	q, err := BindQuery[testQuery](c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Page != 2 {
		t.Errorf("expected page 2, got %d", q.Page)
	}
	if q.PerPage != 20 {
		t.Errorf("expected per_page 20, got %d", q.PerPage)
	}
	if q.Search != "hello" {
		t.Errorf("expected search 'hello', got %q", q.Search)
	}
}

func TestBindQuery_ValidationError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=0", nil)

	_, err := BindQuery[testQuery](c)
	if err == nil {
		t.Fatal("expected validation error for page < 1")
	}
}

func TestMustBindQuery_Success(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=1", nil)

	q := MustBindQuery[testQuery](c)
	if q == nil {
		t.Fatal("expected non-nil query")
	}
	if q.Page != 1 {
		t.Errorf("expected page 1, got %d", q.Page)
	}
}

func TestMustBindQuery_Failure_AbortsWith400(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=-1", nil)

	q := MustBindQuery[testQuery](c)
	if q != nil {
		t.Fatal("expected nil on failure")
	}
	if !c.IsAborted() {
		t.Error("expected context to be aborted")
	}
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- BindURI tests ---

type testURI struct {
	ID int `uri:"id" binding:"required"`
}

func TestBindURI_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	var result *testURI
	var bindErr error
	engine.GET("/items/:id", func(c *gin.Context) {
		result, bindErr = BindURI[testURI](c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/items/42", nil)
	engine.ServeHTTP(w, req)

	if bindErr != nil {
		t.Fatalf("unexpected error: %v", bindErr)
	}
	if result.ID != 42 {
		t.Errorf("expected id 42, got %d", result.ID)
	}
}

func TestMustBindURI_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	var result *testURI
	engine.GET("/items/:id", func(c *gin.Context) {
		result = MustBindURI[testURI](c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/items/7", nil)
	engine.ServeHTTP(w, req)

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.ID != 7 {
		t.Errorf("expected id 7, got %d", result.ID)
	}
}
