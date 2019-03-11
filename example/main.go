package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"net/http"
)

type SampleResource struct {
	Prefix string
}

func (r SampleResource) Get(i int, name string) (gin.H, int) {
	return gin.H{
		"int": i,
		"name": name,
	}, http.StatusOK
}

func main() {
	r := gin.Default()
	gin_restful.Register(r,  SampleResource{Prefix:"/samples"})
	_ = r.Run(":5000")
}
