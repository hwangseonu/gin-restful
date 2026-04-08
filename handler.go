package gin_restful

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func makeHandler(fn func(c *gin.Context) (any, int, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, status, err := fn(c)
		if err != nil {
			handleError(c, err, status)
			return
		}
		if status == http.StatusNoContent {
			c.Status(status)
			return
		}
		c.JSON(status, result)
	}
}

func handleError(c *gin.Context, err error, fallbackStatus int) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		c.AbortWithStatusJSON(httpErr.Status, httpErr)
		return
	}
	status := fallbackStatus
	if status == 0 {
		status = http.StatusInternalServerError
	}
	_ = c.Error(err)
	c.AbortWithStatusJSON(status, gin.H{"message": "internal server error"})
}

func normalizePath(path string) string {
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}
	if len(path) > 1 {
		path = strings.TrimRight(path, "/")
	}
	return path
}
