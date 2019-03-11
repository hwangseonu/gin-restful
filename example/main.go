package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"net/http"
)

type SampleResource struct {
	gin_restful.Resource
}

func (r SampleResource) Get(i int, name string) (gin.H, int) {
	return gin.H{
		"int": i,
		"name": name,
	}, http.StatusNotFound
}

func main() {
	r := gin.Default()
	res := SampleResource{gin_restful.Resource{Prefix:"/samples"}}
	gin_restful.Register(r, res)
	_ = r.Run(":5000")
}
