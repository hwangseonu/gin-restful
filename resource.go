package gin_restful

import (
	"github.com/gin-gonic/gin"
)

type Resource interface {
	RequestBody() interface{}

	Create(body interface{}, c *gin.Context) (gin.H, error)
	Read(id string, c *gin.Context) (gin.H, error)
	ReadAll(c *gin.Context) (gin.H, error)
	Update(id string, body interface{}, c *gin.Context) (gin.H, error)
	Delete(id string, c *gin.Context) (gin.H, error)
}
