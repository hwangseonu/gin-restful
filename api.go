package gin_restful

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerFunc is a function that handles errors returned from resource handlers.
// It receives the gin context, the error, and a fallback status code.
type ErrorHandlerFunc func(c *gin.Context, err error, status int)

// APIOption configures optional settings on the API.
type APIOption func(*API)

// WithErrorHandler sets a custom error handler for the API.
// When set, it replaces the default error handling behavior for all resources
// registered with this API instance.
func WithErrorHandler(handler ErrorHandlerFunc) APIOption {
	return func(api *API) {
		api.errorHandler = handler
	}
}

// API manages RESTful resource registration under a common URL prefix.
type API struct {
	prefix       string
	router       gin.IRouter
	errorHandler ErrorHandlerFunc
}

// NewAPI creates a new API with the given router and URL prefix.
// The router can be a *gin.Engine, *gin.RouterGroup, or any gin.IRouter.
// Optional APIOption arguments can configure error handling and other settings.
func NewAPI(router gin.IRouter, prefix string, opts ...APIOption) *API {
	api := &API{prefix: prefix, router: router}
	for _, opt := range opts {
		opt(api)
	}
	return api
}

// AddResource registers a resource at the given path. The resource is inspected
// via type assertions to determine which HTTP method interfaces it implements.
// Only implemented interfaces become routes. Panics if the resource implements
// none of the handler interfaces (Lister, Getter, Poster, Putter, Patcher, Deleter).
func (api *API) AddResource(path string, resource any) {
	fullPath := normalizePath(api.prefix + path)
	registered := 0

	makeH := func(fn func(c *gin.Context) (any, int, error)) gin.HandlerFunc {
		return makeHandlerWithErrorHandler(fn, api.errorHandler)
	}

	if r, ok := resource.(Poster); ok {
		registered++
		api.router.POST(fullPath, makeH(func(c *gin.Context) (any, int, error) {
			return r.Post(c)
		}))
	}

	if r, ok := resource.(Lister); ok {
		registered++
		api.router.GET(fullPath, makeH(func(c *gin.Context) (any, int, error) {
			return r.List(c)
		}))
	}

	idPath := fullPath + "/:id"

	if r, ok := resource.(Getter); ok {
		registered++
		api.router.GET(idPath, makeH(func(c *gin.Context) (any, int, error) {
			return r.Get(c.Param("id"), c)
		}))
	}

	if r, ok := resource.(Putter); ok {
		registered++
		api.router.PUT(idPath, makeH(func(c *gin.Context) (any, int, error) {
			return r.Put(c.Param("id"), c)
		}))
	}

	if r, ok := resource.(Patcher); ok {
		registered++
		api.router.PATCH(idPath, makeH(func(c *gin.Context) (any, int, error) {
			return r.Patch(c.Param("id"), c)
		}))
	}

	if r, ok := resource.(Deleter); ok {
		registered++
		api.router.DELETE(idPath, makeH(func(c *gin.Context) (any, int, error) {
			return r.Delete(c.Param("id"), c)
		}))
	}

	if registered == 0 {
		panic(fmt.Sprintf("gin-restful: resource at %q implements none of the handler interfaces", path))
	}
}
