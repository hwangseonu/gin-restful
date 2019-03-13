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

//리소스 객체와 URL을 묶어서 사용하기 위한 구조체
type ResourceUrl struct {
	Resource interface{}
	Url      string
}

//리소스들을 관리하고 서버에 등록하기 위한 API구조체
type Api struct {
	App       *gin.RouterGroup
	Prefix    string
	Resources []ResourceUrl
}

//API 구조체의 인스턴스를 인스턴스를 생성하는 함수
func NewApi(app *gin.Engine, prefix string) *Api {
	return &Api{
		App:       app.Group(prefix),
		Prefix:    prefix,
		Resources: make([]ResourceUrl, 0),
	}
}

//API 인스턴스에 새로운 Resource를 등록하는 메서드
func (a *Api) AddResource(resource interface{}, url string) {
	if a.App != nil {
		a.registerResource(resource, url)
	} else {
		a.Resources = append(a.Resources, ResourceUrl{
			Resource: resource,
			Url: url,
		})
	}
}

//API 인스턴스에 등록 된 리소스 핸들러들을 gin.HandlerChain 타입으로 반환하는 메서드
func (a *Api) GetHandlersChain() gin.HandlersChain {
	result := make([]gin.HandlerFunc, 0)
	for _, v := range a.Resources {
		resource := v.Resource
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

//API 에 등록된 리소스를 gin 서버에 등록시키는 메서드
func (a *Api) registerResource(resource interface{}, url string) {
	for i := 0; i < reflect.TypeOf(resource).NumMethod(); i++ {
		value := reflect.ValueOf(resource)
		method := reflect.TypeOf(resource).Method(i)
		if !isHttpMethod(method.Name) {
			continue
		}
		args := parseArgs(method)
		url := createUrl(url, args)
		g := a.App.Group(url)
		for _, m := range parseMiddlewares(resource, method.Name) {
			g.Use(m)
		}
		g.Handle(strings.ToUpper(method.Name), "", createHandlerFunc(value, method, args))
	}
}

//인자로 받은 메서드 이름이 http 메서드인지 확인하는 함수
func isHttpMethod(name string) bool {
	for _, k := range httpmethods {
		if strings.ToUpper(name) == k {
			return true
		}
	}
	return false
}

//리소스에 등록되어 있는 메서드의 인자들을 http url 로 만들어 반환하는 함수
func createUrl(url string, args []string) string {
	for i, a := range args {
		if a == "context" {
			continue
		}
		url += "/:" + a + strconv.Itoa(i)
	}
	return url
}

//메서드의 인자 타입을 반환하는 함수
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

//인자 타입을 바탕으로 리소스의 메서드를 실행시키기 위한 parameter 를 만드는 함수
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
					Message: "argument " + arg + strconv.Itoa(i) + " is must int",
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

//gin 서버에 등록 가능한 handler 함수를 만들어 반환하는 함수
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

//리소스 메서드에 등록된 middleware 들을 반환하는 함수
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