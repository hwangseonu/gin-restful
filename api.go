package gin_restful

import (
	"github.com/gin-gonic/gin"
)

type API struct {
	Prefix    string
	resources map[string]Resource
}

func NewAPI(prefix string) *API {
	api := new(API)
	api.Prefix = prefix
	api.resources = make(map[string]Resource)
	return api
}

func (api *API) RegisterResource(path string, resource Resource) {
	api.resources[path] = resource
}

func (api *API) RegisterHandlers(router *gin.RouterGroup) {
	for path, resource := range api.resources {
		path = api.Prefix + path
		handler := gin.HandlerFunc(func(c *gin.Context) {
			handleHTTP(resource, c)
		})

		if resource.Create != nil {
			router.POST(path, handler)
		}

		if resource.ReadAll != nil {
			router.GET(path, handler)
		}

		path = path + "/:id"

		if resource.Read != nil {
			router.GET(path, handler)
		}

		if resource.Update != nil {
			router.PUT(path, handler)
			router.PATCH(path, handler)
		}

		if resource.Delete != nil {
			router.DELETE(path, handler)
		}
	}
}
