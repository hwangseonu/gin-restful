package gin_restful

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

func JsonRequired(c *gin.Context, json interface{}) (interface{}, error) {
	mustType := reflect.TypeOf(json)
	value := reflect.New(mustType)
	if err := c.ShouldBindJSON(value.Interface()); err != nil {
		return nil, err
	}
	return value.Elem().Interface(), nil
}
