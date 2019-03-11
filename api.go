package gin_restful

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Api struct {
	*gin.Engine
	Prefix string
}

func NewApi(e *gin.Engine, prefix string) *Api {
	return &Api{
		Engine: e,
		Prefix: prefix,
	}
}

func (a *Api) AddResource(resource interface{}, url string) {
	for i := 0; i < reflect.TypeOf(resource).NumMethod(); i++ {
		value := reflect.ValueOf(resource)
		method := reflect.TypeOf(resource).Method(i)
		url, args := makeUrl(url, method)
		a.Engine.Handle(strings.ToUpper(method.Name), a.Prefix + url, makeFunc(value, method, args))
	}
}

func makeUrl(url string, method reflect.Method) (string, []string) {
	args := make([]string, 0)
	for i := 1; i < method.Type.NumIn(); i++ {
		arg := method.Type.In(i).String()
		args = append(args, arg)
		if method.Type.In(i).String() == "*gin.Context" {
			arg = "context"
			continue
		}
		url += "/:" + arg + strconv.Itoa(i)
	}
	return url, args
}

func makeValues(c *gin.Context, resource reflect.Value, args []string) ([]reflect.Value, error){
	values := []reflect.Value{resource}
	for i, arg := range args {
		p := c.Param(arg+strconv.Itoa(i+1))
		switch arg {
		case "string":
			values = append(values, reflect.ValueOf(p))
			break
		case "int":
			if num, err := strconv.Atoi(p); err != nil {
				return []reflect.Value{}, ApplicationError{
					Message: "argument " + arg + strconv.Itoa(i) + "is must int",
					Status:  http.StatusBadRequest,
				}
			} else {
				values = append(values, reflect.ValueOf(num))
			}
			break
		case "float64":
			if num, err := strconv.ParseFloat(p, 64); err != nil {
				return []reflect.Value{}, ApplicationError{
					Message: "argument " + arg + strconv.Itoa(i) + "is must " + arg,
					Status:  http.StatusBadRequest,
				}
			} else {
				values = append(values, reflect.ValueOf(num))
			}
			break
		case "bool":
			if p == "" || p == "false" || p == "0" || p == "null" || p == "nil" || p == "off" {
				values = append(values, reflect.ValueOf(false))
			} else {
				values = append(values, reflect.ValueOf(true))
			}
			break
		case "context":
			values = append(values, reflect.ValueOf(c))
			break
		default:
			panic(errors.New("method argument can string, integer"))
		}
	}
	return values, nil
}

func makeFunc(resource reflect.Value, method reflect.Method, args []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		values, err := makeValues(c, resource, args)
		if err != nil {
			ae, ok := err.(ApplicationError)
			status := http.StatusInternalServerError
			if ok {
				status = ae.Status
			}
			c.JSON(status, gin.H{"message": err.Error()})
			return
		}
		returns := method.Func.Call(values)
		status := http.StatusOK
		switch len(returns) {
		case 0:
			c.Status(http.StatusOK)
			return
		case 2:
			status = int(returns[1].Int())
		}
		if _, err := json.MarshalIndent(returns[0].Interface(), "", "  "); err != nil {
			c.String(status, returns[0].String())
		} else {
			c.JSON(status, returns[0].Interface())
		}
	}
}
