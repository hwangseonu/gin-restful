package gin_restful

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type API struct {
	prefix string
	router gin.IRouter
}

func NewAPI(router gin.IRouter, prefix string) *API {
	return &API{prefix: prefix, router: router}
}

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
