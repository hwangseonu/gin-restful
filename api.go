package gin_restful

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// API manages RESTful resource registration under a common URL prefix.
type API struct {
	prefix string
	router gin.IRouter
}

// NewAPI creates a new API with the given router and URL prefix.
// The router can be a *gin.Engine, *gin.RouterGroup, or any gin.IRouter.
func NewAPI(router gin.IRouter, prefix string) *API {
	return &API{prefix: prefix, router: router}
}

// AddResource registers a resource at the given path. The resource is inspected
// via type assertions to determine which HTTP method interfaces it implements.
// Only implemented interfaces become routes. Panics if the resource implements
// none of the handler interfaces (Lister, Getter, Poster, Putter, Patcher, Deleter).
func (api *API) AddResource(path string, resource any) {
	fullPath := normalizePath(api.prefix + path)
	registered := 0

	if r, ok := resource.(Poster); ok {
		registered++
		api.router.POST(fullPath, makeHandler(func(c *gin.Context) (any, int, error) {
			return r.Post(c)
		}))
	}

	if r, ok := resource.(Lister); ok {
		registered++
		api.router.GET(fullPath, makeHandler(func(c *gin.Context) (any, int, error) {
			return r.List(c)
		}))
	}

	idPath := fullPath + "/:id"

	if r, ok := resource.(Getter); ok {
		registered++
		api.router.GET(idPath, makeHandler(func(c *gin.Context) (any, int, error) {
			return r.Get(c.Param("id"), c)
		}))
	}

	if r, ok := resource.(Putter); ok {
		registered++
		api.router.PUT(idPath, makeHandler(func(c *gin.Context) (any, int, error) {
			return r.Put(c.Param("id"), c)
		}))
	}

	if r, ok := resource.(Patcher); ok {
		registered++
		api.router.PATCH(idPath, makeHandler(func(c *gin.Context) (any, int, error) {
			return r.Patch(c.Param("id"), c)
		}))
	}

	if r, ok := resource.(Deleter); ok {
		registered++
		api.router.DELETE(idPath, makeHandler(func(c *gin.Context) (any, int, error) {
			return r.Delete(c.Param("id"), c)
		}))
	}

	if registered == 0 {
		panic(fmt.Sprintf("gin-restful: resource at %q implements none of the handler interfaces", path))
	}
}
