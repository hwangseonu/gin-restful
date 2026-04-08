package restful

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

func benchmarkEngine() *gin.Engine {
	engine := gin.New()
	api := NewAPI(engine, "/api")
	api.AddResource("/items", &fullCRUDResource{})
	return engine
}

func BenchmarkHandler_List(b *testing.B) {
	engine := benchmarkEngine()
	req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
	for b.Loop() {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

func BenchmarkHandler_Get(b *testing.B) {
	engine := benchmarkEngine()
	req := httptest.NewRequest(http.MethodGet, "/api/items/1", nil)
	for b.Loop() {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

func BenchmarkHandler_Post(b *testing.B) {
	engine := benchmarkEngine()
	body := []byte(`{"name":"test","age":25}`)
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

func BenchmarkHandler_Delete(b *testing.B) {
	engine := benchmarkEngine()
	req := httptest.NewRequest(http.MethodDelete, "/api/items/1", nil)
	for b.Loop() {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

func BenchmarkHandler_ErrorPath(b *testing.B) {
	engine := gin.New()
	api := NewAPI(engine, "/api")
	api.AddResource("/items", &errorResource{})
	req := httptest.NewRequest(http.MethodGet, "/api/items/1", nil)
	for b.Loop() {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}
