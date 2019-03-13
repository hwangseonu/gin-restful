package gin_restful

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type Resource struct {
	Middlewares map[string][]gin.HandlerFunc
}

func InitResource() *Resource {
	return &Resource{
		Middlewares: make(map[string][]gin.HandlerFunc, 0),
	}
}

func (r *Resource) AddMiddleware(middleware gin.HandlerFunc, methods ...string) {
	for _, m := range methods {
		m = strings.ToUpper(m)
		middlewares := r.Middlewares[m]
		middlewares = append(middlewares, middleware)
		r.Middlewares[m] = middlewares
	}
}
