package gin_restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Bind[T any](c *gin.Context) (*T, error) {
	var body T
	if err := c.ShouldBind(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

func MustBind[T any](c *gin.Context) *T {
	body, err := Bind[T](c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &HTTPError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return nil
	}
	return body
}
