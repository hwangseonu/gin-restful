package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"net/http"
)

type SampleResource struct {
	*gin_restful.Resource
}

func (r SampleResource) Get(name string) (gin.H, int) {
	return gin.H{
		"name": name,
	}, http.StatusOK
}

func (r SampleResource) Post(name string) string {
	return name
}

func SampleMiddleware(c *gin.Context) {
	println("Hello, World")
}

func main() {
	r := gin.Default()
	v1 := gin_restful.NewApi(r, "/")
	res := SampleResource{gin_restful.InitResource()}
	res.AddMiddleware(SampleMiddleware, http.MethodGet)
	v1.AddResource(res, "/samples")
	_ = r.Run(":5000")
}
