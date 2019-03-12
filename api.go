package gin_restful

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var httpmethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodPut,
	http.MethodHead,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

type ResourceTuple struct {
	Resource interface{}
	Url      string
}

type Api struct {
	App       *gin.RouterGroup
	Prefix    string
	Resources []ResourceTuple
}

func NewApi(app *gin.Engine, prefix string) *Api {
	return &Api{
		App:       app.Group(prefix),
		Prefix:    prefix,
		Resources: make([]ResourceTuple, 0),
	}
}

func (a *Api) AddResource(resource interface{}, url string) {
	if a.App != nil {
		a.registerResource(resource, url)
	} else {
		a.Resources = append(a.Resources, ResourceTuple{
			Resource: resource,
			Url: url,
		})
	}
}

func (a *Api) GetHandlersChain() gin.HandlersChain {
	result := make([]gin.HandlerFunc, 0)
	for _, tuple := range a.Resources {
		resource := tuple.Resource
		for i := 0; i < reflect.TypeOf(resource).NumMethod(); i++ {
			value := reflect.ValueOf(resource)
			method := reflect.TypeOf(resource).Method(i)
			if !isHttpMethod(method.Name) {
				continue
			}
			args := parseArgs(method)
			result = append(result, createHandlerFunc(value, method, args))
		}
	}
	return result
}

func (a *Api) registerResource(resource interface{}, url string) {
	for i := 0; i < reflect.TypeOf(resource).NumMethod(); i++ {
		value := reflect.ValueOf(resource)
		method := reflect.TypeOf(resource).Method(i)
		if !isHttpMethod(method.Name) {
			continue
		}
		args := parseArgs(method)
		url := createUrl(url, args)
		a.App.Handle(strings.ToUpper(method.Name), url, createHandlerFunc(value, method, args))
	}
}

func isHttpMethod(name string) bool {
	for _, k := range httpmethods {
		if strings.ToUpper(name) == k {
			return true
		}
	}
	return false
}

func createUrl(url string, args []string) string {
	for i, a := range args {
		if a == "context" {
			continue
		}
		url += "/:" + a + strconv.Itoa(i)
	}
	return url
}

func parseArgs(method reflect.Method) []string {
	args := make([]string, 0)
	for i := 1; i < method.Type.NumIn(); i++ {
		arg := method.Type.In(i).String()
		if arg == "*gin.Context" {
			arg = "context"
		}
		args = append(args, arg)
	}
	return args
}

func createValues(c *gin.Context, resource reflect.Value, args []string) ([]reflect.Value, error) {
	values := []reflect.Value{resource}
	for i, arg := range args {
		p := c.Param(arg + strconv.Itoa(i))
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

func createHandlerFunc(resource reflect.Value, method reflect.Method, args []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		values, err := createValues(c, resource, args)
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
