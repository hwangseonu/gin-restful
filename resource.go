package gin_restful

import (
	"log"
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
	id := c.Param("id")
	body := resource.RequestBody(c.Request.Method)

	v := reflect.ValueOf(body)

	if v.IsValid() && !v.IsNil() {
		if err := c.ShouldBindJSON(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
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
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}
		result, status, err = resource.Update(id, body, c)
	case "DELETE":
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}
		result, status, err = resource.Delete(id, c)
	}

	if err != nil {
		err = c.Error(err)
		c.JSON(status, gin.H{"error": err})
		log.Println(err)
	} else {
		c.JSON(status, result)
	}
}
