package gin_restful

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
)

func Register(e *gin.Engine, resource interface{}) {
	registerGet(e, resource)
}

func registerGet(e *gin.Engine, resource interface{}) {
	get, ok := reflect.TypeOf(resource).MethodByName("Get")
	if ok {
		url := reflect.ValueOf(resource).FieldByName("Prefix").String()
		args := make([]string, 0)
		for i := 1; i < get.Type.NumIn(); i++ {
			arg := get.Type.In(i)
			args = append(args, arg.String())
			url += "/:"+arg.String()+strconv.Itoa(i)
		}
		e.GET(url, func (c *gin.Context) {
			values := make([]reflect.Value, 0)
			values = append(values, reflect.ValueOf(resource))
			for i, v := range args  {
				p := c.Param(v + strconv.Itoa(i+1))
				if v == "string" {
					values = append(values, reflect.ValueOf(p))
				} else if v == "int" {
					if j, err := strconv.Atoi(p); err != nil {
						c.AbortWithStatus(http.StatusBadRequest)
						return
					} else {
						values = append(values, reflect.ValueOf(j))
					}
				}
			}
			get.Func.Call(values)
			c.Status(http.StatusOK)
		})
	}
}