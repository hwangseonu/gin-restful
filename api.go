//gin_restful 은 gin 을 이용한 restful api 를 간편하게 만들기 위한 extension 입니다.
//go 언어로 restful api 를 더 편하게 만들고 싶어 개발하였습니다.
package gin_restful

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
)

//Api 구조체는 Resource 인스턴스들을 관리하고 gin 서버에 등록하기 위한 구조체입니다.
//NewApi 함수로 인스턴스를 생성하여 사용합니다.
//AddResource 함수로 Api 인스턴스에 Resource 를 등록할 수 있습니다.
type Api struct {
	App       *gin.RouterGroup
	Prefix    string
	Resources map[string]interface{}
}

//Api 구조체의 인스턴스를 인스턴스를 생성하여 포인터로 반환하는 함수입니다.
//첫번째 인자 app(type: *gin.Engine)은 Resource 를 등록 할 gin 서버의 인스턴스입니다.
//두번째 인자 prefix(type string)은 api url 의 제일 앞 부분에 붙습니다.
func NewApi(app *gin.Engine, prefix string) *Api {
	return &Api{
		App:       app.Group(prefix),
		Prefix:    prefix,
		Resources: make(map[string]interface{}),
	}
}

//Api 인스턴스에 새로운 Resource 를 등록하는 메서드입니다.
//Api 의 필드 App 이 nil 이 아니라면 gin 서버에 즉시 등록하고 nil 이라면 Api 구조체에 잠시 저장합니다.
//서버에 등록할 때는 Resource 의 메서드 중 http method 인 것을 찾아 인자를 파싱하여 url 을 생성합니다.
//요청을 받았을때 메서드가 실행되며 각각의 인자는 자동으로 채워집니다.
//string, int, float, bool 타입의 인자는 url 에서 파싱하여 전달합니다.
//*gin.Context 타입의 인자는 해당 요청의 context 로 채워집니다.
//구조체 타입의 인자는 하나만 존재할 수 있으며 요청의 body 를 파싱하여 채워집니다.
func (a *Api) AddResource(resource interface{}, url string) {
	if a.App != nil {
		a.registerResource(resource, url)
	} else {
		a.Resources[url] = resource
	}
}

//Api 인스턴스에 등록된 Resource 의 Handler 들을 gin.HandlerChain 타입으로 반환하는 메서드입니다.
func (a *Api) GetHandlersChain() gin.HandlersChain {
	result := make([]gin.HandlerFunc, 0)
	for _, v := range a.Resources {
		for i := 0; i < reflect.TypeOf(v).NumMethod(); i++ {
			value := reflect.ValueOf(v)
			method := reflect.TypeOf(v).Method(i)
			if !isHttpMethod(method.Name) {
				continue
			}
			args := parseArgs(method)
			result = append(result, createHandlerFunc(value, method, args))
		}
	}
	return result
}

//새로운 gin 서버를 Api 에 등록하고 Api 에 저장되어있던 Resource 들을 gin 서버에 등록시키는 메서드입니다.
func (a *Api) InitApp(e *gin.Engine) {
	a.App = e.Group(a.Prefix)
	for k, v := range a.Resources {
		a.registerResource(v, k)
	}
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
		g := a.App.Group(url, parseMiddlewares(resource, method.Name)...)
		g.Handle(strings.ToUpper(method.Name), "", createHandlerFunc(value, method, args))
	}
}