package restful_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	restful "github.com/hwangseonu/gin-restful"
)

type ExampleItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ExampleResource struct{}

func (r *ExampleResource) List(c *gin.Context) (any, int, error) {
	return []ExampleItem{{ID: "1", Name: "item1"}}, http.StatusOK, nil
}

func (r *ExampleResource) Get(id string, c *gin.Context) (any, int, error) {
	if id == "1" {
		return ExampleItem{ID: "1", Name: "item1"}, http.StatusOK, nil
	}
	return nil, 0, restful.Abort(http.StatusNotFound, "not found")
}

func ExampleNewAPI() {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := restful.NewAPI(engine, "/api/v1")
	api.AddResource("/items", &ExampleResource{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/items", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output:
	// 200
}

type CreateReq struct {
	Name string `json:"name" binding:"required"`
}

type ExamplePostResource struct{}

func (r *ExamplePostResource) Post(c *gin.Context) (any, int, error) {
	body, err := restful.Bind[CreateReq](c)
	if err != nil {
		return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
	}
	return gin.H{"name": body.Name}, http.StatusCreated, nil
}

func ExampleBind() {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := restful.NewAPI(engine, "/api")
	api.AddResource("/items", &ExamplePostResource{})

	req := httptest.NewRequest(http.MethodPost, "/api/items", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output:
	// 400
}

func ExampleAbort() {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api := restful.NewAPI(engine, "/api")
	api.AddResource("/items", &ExampleResource{})

	req := httptest.NewRequest(http.MethodGet, "/api/items/999", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	fmt.Println(w.Code)
	// Output:
	// 404
}

func ExampleAbort_withCodeAndDetails() {
	err := restful.Abort(422, "validation failed",
		restful.WithCode("VALIDATION_ERROR"),
		restful.WithDetails(map[string]string{"name": "required"}),
	)
	fmt.Println(err.Status, err.Code, err.Message)
	// Output:
	// 422 VALIDATION_ERROR validation failed
}
