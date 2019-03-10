package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
)

type BaseResource struct {
	gin_restful.Resource
}

func (r *BaseResource) Get(name int, a string) {
	println(name)
	println(a)
}

func main() {
	r := gin.Default()
	res := BaseResource{gin_restful.Resource{Prefix:"/bases"}}
	gin_restful.Register(r, res)
	_ = r.Run(":5000")
}
