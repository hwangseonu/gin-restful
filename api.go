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

func (api *API) RegisterHandlers(engine *gin.Engine) {
	for path, resource := range api.resources {
		path = api.Prefix + "/" + path
		engine.Any(path, func(context *gin.Context) {
			handleHTTP(resource, context)
		})
		engine.Any(path+"/:id", func(context *gin.Context) {
			handleHTTP(resource, context)
		})
	}
}
