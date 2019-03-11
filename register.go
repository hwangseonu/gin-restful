package gin_restful

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func Register(e *gin.Engine, resource interface{}) {
	numMethod := reflect.TypeOf(resource).NumMethod()
	for i := 0 ; i < numMethod; i++ {
		v := reflect.ValueOf(resource)
		method := reflect.TypeOf(resource).Method(i)
		url, args := makeUrl(v.FieldByName("Prefix").String(), method)
		e.Handle(strings.ToUpper(method.Name), url, makeFunc(v, method, args))
	}
}

func makeUrl(prefix string, method reflect.Method) (string, []string) {
	url := prefix
	args := make([]string, 0)
	for i := 1; i < method.Type.NumIn(); i++ {
		arg := method.Type.In(i).String()
		args = append(args, arg)
		if arg == "*gin.Context" {
			continue
		}
		url += "/:"+arg+strconv.Itoa(i)
	}
	return url, args
}

func makeValues(params gin.Params, v reflect.Value, args []string, c *gin.Context) ([]reflect.Value, error) {
	values := []reflect.Value{v}
	for i, v := range args {
		p := params.ByName(v + strconv.Itoa(i+1))
		if v == "string" {
			values = append(values, reflect.ValueOf(p))
		} else if v == "int" {
			if j, err := strconv.Atoi(p); err != nil {
				return []reflect.Value{}, ApplicationError{
					Message: v+strconv.Itoa(i+1) + " is must int",
					Status: http.StatusBadRequest,
				}
			} else {
				values = append(values, reflect.ValueOf(j))
			}
		} else if v == "bool" {
			if p == "false" || p == "0" || p == "off" || p == "null" {
				values = append(values, reflect.ValueOf(false))
			} else {
				values = append(values, reflect.ValueOf(true))
			}
		} else if v == "*gin.Context" {
			values = append(values, reflect.ValueOf(c))
		}
	}
	return values, nil
}

func makeFunc(v reflect.Value, method reflect.Method, args []string) func (c *gin.Context){
	return func(c *gin.Context) {
		values, err := makeValues(c.Params, v, args, c)
		if err != nil {
			ae, ok := err.(ApplicationError)
			status := http.StatusInternalServerError
			if ok {
				status = ae.Status
			}
			c.AbortWithStatusJSON(status, gin.H{"message": err.Error()})
			return
		}
		returns := method.Func.Call(values)
		status := http.StatusOK
		if len(returns) == 0 {
			c.Status(http.StatusOK)
			return
		}
		if len(returns) == 2 {
			status = int(returns[1].Int())
		}
		if _, err := json.Marshal(returns[0].Interface()); err != nil {
			c.String(status, returns[0].String())
		} else {
			c.JSON(status, returns[0].Interface())
		}
	}
}