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

func getJsonData(c *gin.Context, json interface{}) (interface{}, error) {
	mustType := reflect.TypeOf(json)
	value := reflect.New(mustType)
	if err := c.ShouldBindJSON(value.Interface()); err != nil {
		return nil, err
	}
	return value.Elem().Interface(), nil
}

func isHttpMethod(name string) bool {
	for _, k := range httpmethods {
		if strings.ToUpper(name) == k {
			return true
		}
	}
	return false
}

func createUrl(url string, args []reflect.Type) string {
	for i, a := range args {
		if a.Kind() == reflect.Struct {
			continue
		}
		if a.String() == "*gin.Context" {
			continue
		}
		url += "/:" + a.String() + strconv.Itoa(i)
	}
	return url
}

func parseArgs(method reflect.Method) []reflect.Type {
	args := make([]reflect.Type, 0)
	can := []string{"string", "int", "float64", "bool", "*gin.Context"}

	addedStruct := false
	for i := 1; i < method.Type.NumIn(); i++ {
		arg := method.Type.In(i)
		if arg.Kind() == reflect.Struct {
			if addedStruct {
				panic(errors.New("method argument can string, int, float64, bool, *gin.Context, one struct"))
			} else {
				addedStruct = true
				args = append(args, arg)
				continue
			}
		}
		if !contains(can, arg.String()) {
			panic(errors.New("method argument can string, int, float64, bool, *gin.Context, one struct"))
		}
		args = append(args, arg)
	}
	return args
}

func createValues(c *gin.Context, resource reflect.Value, args []reflect.Type) ([]reflect.Value, error) {
	values := []reflect.Value{resource}

	for i, arg := range args {
		p := c.Param(arg.String() + strconv.Itoa(i))
		switch arg.Kind() {
		case reflect.String:
			values = append(values, reflect.ValueOf(p))
			break
		case reflect.Int:
			if num, err := strconv.Atoi(p); err != nil {
				return []reflect.Value{}, ApplicationError{
					Message: "argument " + arg.String() + strconv.Itoa(i) + " is must int",
					Status:  http.StatusBadRequest,
				}
			} else {
				values = append(values, reflect.ValueOf(num))
			}
			break
		case reflect.Float64:
			if num, err := strconv.ParseFloat(p, 64); err != nil {
				return []reflect.Value{}, ApplicationError{
					Message: "argument " + arg.String() + strconv.Itoa(i) + "is must float64",
					Status:  http.StatusBadRequest,
				}
			} else {
				values = append(values, reflect.ValueOf(num))
			}
			break
		case reflect.Bool:
			if p == "" || p == "false" || p == "0" || p == "null" || p == "nil" || p == "off" {
				values = append(values, reflect.ValueOf(false))
			} else {
				values = append(values, reflect.ValueOf(true))
			}
			break
		case reflect.Ptr:
			if arg.String() == "*gin.Context" {
				values = append(values, reflect.ValueOf(c))
			}
			break
		case reflect.Struct:
			body := reflect.New(arg).Elem().Interface()
			if v, err := getJsonData(c, body); err != nil {
				return []reflect.Value{}, ApplicationError{
					Message: err.Error(),
					Status:  400,
				}
			} else {
				values = append(values, reflect.ValueOf(v))
			}
			break
		default:
			panic(errors.New("method argument can string, int, float64, bool, *gin.Context, one struct"))
		}
	}
	return values, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func createHandlerFunc(resource reflect.Value, method reflect.Method, args []reflect.Type) gin.HandlerFunc {
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

func parseMiddlewares(resource interface{}, method string) []gin.HandlerFunc {
	r := reflect.ValueOf(resource)
	method = strings.ToUpper(method)
	f := r.FieldByName("Middlewares")
	if !f.IsValid() {
		return make([]gin.HandlerFunc, 0)
	}
	middlewares := r.FieldByName("Middlewares").Interface().(map[string][]gin.HandlerFunc)
	return middlewares[method]
}
