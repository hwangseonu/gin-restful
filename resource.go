package gin_restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resource interface {
	RequestBody() interface{}

	Create(body interface{}, c *gin.Context) (gin.H, int, error)
	Read(id string, c *gin.Context) (gin.H, int, error)
	ReadAll(c *gin.Context) (gin.H, int, error)
	Update(id string, body interface{}, c *gin.Context) (gin.H, int, error)
	Delete(id string, c *gin.Context) (gin.H, int, error)
}

func handleHTTP(resource Resource, c *gin.Context) {
	id := c.Param("id")
	body := resource.RequestBody()

	if err := c.ShouldBindBodyWithJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result gin.H
	var status int
	var err error

	switch c.Request.Method {
	case "POST":
		if id == "" {
			result, status, err = resource.Create(body, c)
		} else {
			result, status, err = resource.Update(id, body, c)
		}
	case "GET":
		if id == "" {
			result, status, err = resource.ReadAll(c)
		} else {
			result, status, err = resource.Read(id, c)
		}
	case "PUT", "PATCH":
		result, status, err = resource.Update(id, body, c)
	case "DELETE":
		result, status, err = resource.Delete(id, c)
	}

	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
	} else {
		c.JSON(status, result)
	}
}
