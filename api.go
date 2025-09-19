package gin_restful

import (
	"net/http"

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
		//engine.Handle("GET", path, makeHandler("GET", resource))
		engine.Handle("POST", path, makeHandler("POST", resource))

	}

}

func makeHandler(_ string, r Resource) gin.HandlerFunc {
	return func(c *gin.Context) {
		body := r.RequestBody()
		if err := c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		result, err := r.Create(body, c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusCreated, result)
	}
}
