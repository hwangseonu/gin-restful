package gin_restful

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Resource interface {
	RequestBody(method string) any

	Create(body interface{}, c *gin.Context) (gin.H, int, error)
	Read(id string, c *gin.Context) (gin.H, int, error)
	ReadAll(c *gin.Context) (gin.H, int, error)
	Update(id string, body interface{}, c *gin.Context) (gin.H, int, error)
	Delete(id string, c *gin.Context) (gin.H, int, error)
}

func handleHTTP(resource Resource, c *gin.Context) {
	var result gin.H
	var status int
	var err error

	id := c.Param("id")
	body := resource.RequestBody(c.Request.Method)
	v := reflect.ValueOf(body)

	if v.IsValid() && !v.IsNil() {
		if err := c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if id == "" {
		switch c.Request.Method {
		case "POST":
			result, status, err = resource.Create(body, c)
		case "GET":
			result, status, err = resource.ReadAll(c)
		}
	} else {
		switch c.Request.Method {
		case "GET":
			result, status, err = resource.Read(id, c)
		case "PUT", "PATCH":
			result, status, err = resource.Update(id, body, c)
		case "DELETE":
			result, status, err = resource.Delete(id, c)
		}
	}

	if err != nil {
		_ = c.Error(err)
		c.JSON(status, gin.H{"error": err})
	} else {
		c.JSON(status, result)
	}
}
