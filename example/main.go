package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"net/http"
)

type SampleResource struct {
	Prefix string
}

func (r SampleResource) Get() (gin.H, int) {
	return "Hello, World", http.StatusOK
}

func main() {
	r := gin.Default()
	res := SampleResource{Prefix:"/samples"}
	gin_restful.Register(r, res)
	_ = r.Run(":5000")
}
