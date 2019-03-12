package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"net/http"
)

type SampleResource struct {
}

func (r SampleResource) Get(name string) (gin.H, int) {
	return gin.H{
		"name": name,
	}, http.StatusOK
}

func main() {
	r := gin.Default()
	v1 := gin_restful.NewApi(r, "/")
	v1.AddResource(SampleResource{}, "/samples")
	_ = r.Run(":5000")
}
