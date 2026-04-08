package gin_restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Bind binds the request body to T using gin's ShouldBind.
// The binding method is determined by the Content-Type header.
func Bind[T any](c *gin.Context) (*T, error) {
	var body T
	if err := c.ShouldBind(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

// MustBind binds the request body to T and returns a pointer to it.
// On failure, it automatically responds with 400 Bad Request and aborts the
// middleware chain. Callers must check for a nil return and exit early:
//
//	body := restful.MustBind[MyReq](c)
//	if body == nil {
//		return nil, 0, nil // already aborted
//	}
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

// BindQuery binds URL query parameters to T using gin's ShouldBindQuery.
func BindQuery[T any](c *gin.Context) (*T, error) {
	var query T
	if err := c.ShouldBindQuery(&query); err != nil {
		return nil, err
	}
	return &query, nil
}

// MustBindQuery binds URL query parameters to T. On failure, it automatically
// responds with 400 Bad Request and aborts the middleware chain.
func MustBindQuery[T any](c *gin.Context) *T {
	query, err := BindQuery[T](c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &HTTPError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return nil
	}
	return query
}

// BindURI binds URI path parameters to T using gin's ShouldBindUri.
func BindURI[T any](c *gin.Context) (*T, error) {
	var params T
	if err := c.ShouldBindUri(&params); err != nil {
		return nil, err
	}
	return &params, nil
}

// MustBindURI binds URI path parameters to T. On failure, it automatically
// responds with 400 Bad Request and aborts the middleware chain.
func MustBindURI[T any](c *gin.Context) *T {
	params, err := BindURI[T](c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &HTTPError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return nil
	}
	return params
}
