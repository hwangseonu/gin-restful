//gin_restful 은 gin 을 이용한 restful api 를 간편하게 만들기 위한 extension 입니다.
//go 언어로 restful api 를 더 편하게 만들고 싶어 개발하였습니다.
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

//ResourceUrl 은 Resource 의 인스턴스와 URL 을 묶어서 사용하기 위한 구조체입니다.
type ResourceUrl struct {
	Resource interface{}
	Url      string
}

//Api 구조체는 Resource 인스턴스들을 관리하고 gin 서버에 등록하기 위한 구조체입니다.
//NewApi 함수로 인스턴스를 생성하여 사용합니다.
//AddResource 함수로 Api 인스턴스에 Resource 를 등록할 수 있습니다.
type Api struct {
	App       *gin.RouterGroup
	Prefix    string
	Resources []ResourceUrl
}

//Api 구조체의 인스턴스를 인스턴스를 생성하여 포인터로 반환하는 함수입니다.
//첫번째 인자 app(type: *gin.Engine)은 Resource 를 등록 할 gin 서버의 인스턴스입니다.
//두번째 인자 prefix(type string)은 api url 의 제일 앞 부분에 붙습니다.
func NewApi(app *gin.Engine, prefix string) *Api {
	return &Api{
		App:       app.Group(prefix),
		Prefix:    prefix,
		Resources: make([]ResourceUrl, 0),
	}
}

//Api 인스턴스에 새로운 Resource 를 등록하는 메서드입니다.
//Api 의 필드 App 이 nil 이 아니라면 gin 서버에 즉시 등록하고 nil 이라면 Api 구조체에 잠시 저장합니다.
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

//Api 인스턴스에 등록된 Resource 의 Handler 들을 gin.HandlerChain 타입으로 반환하는 메서드입니다.
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

//Api 인스턴스에 등록된 Resource 를 gin 서버에 등록시키는 메서드입니다.
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

//인자로 받은 메서드 이름이 http 메서드인지 확인하여 bool 로 반환하는 함수입니다.
func isHttpMethod(name string) bool {
	for _, k := range httpmethods {
		if strings.ToUpper(name) == k {
			return true
		}
	}
	return false
}

//Resource 에 등록되어 있는 Handler 의 인자들을 http url 로 파싱하여 반환하는 함수입니다.
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

//메서드의 인자 타입들을 슬라이스로 반환하는 함수입니다.
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

//args(type: []string) 으로 Resource 의 Handler 를 실행시키기 위한 parameter 들을 만드는 함수입니다.
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
			if v, err := JsonRequired(c, body); err != nil {
				return []reflect.Value{}, ApplicationError{
					Message: err.Error(),
					Status: 400,
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

//gin 서버에 등록 가능한 형태 handler 를 만들어 반환하는 함수입니다.
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

//Resource 인스턴스에 등록된 middleware 들을 반환하는 함수입니다.
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