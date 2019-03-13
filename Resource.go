package gin_restful

import (
	"github.com/gin-gonic/gin"
	"strings"
)

/*
	Resource 는 기본적으로 Api 인스턴스에 등록 가능한 Resource 의 형태입니다.
	각 Handler 에 적용될 Middleware 들을 포함합니다.
	Resource 포인터를 임베딩하여 사용할 수 있지만 Middleware 가 필요없다면 임베딩 하지 않아도 무방합니다.
	Resource 포인터를 임베딩할 때 InitResource 함수를 사용합니다.
	Resource 에 Handler 를 등록할 때는 사용할 http 메서드와 같은 메서드를 정의하면 됩니다.
	Handler 의 이름으로는 Get, Post, Put, Patch, Delete 등을 사용할 수 있습니다.
	Handler 이름의 첫글자는 무조건 대문자여야합니다.
	Handler 의 에는 *gin.Context, string, int, bool, float64 등을 사용할 수 있으며 *gin.Context 를 제외한 인자는 path variable 이 됩니다.
*/
type Resource struct {
	Middlewares map[string][]gin.HandlerFunc
}

//Resource 구조체를 초기화하여 포인터로 반환해주는 함수입니다.
func InitResource() *Resource {
	return &Resource{
		Middlewares: make(map[string][]gin.HandlerFunc, 0),
	}
}

//Resource 인스턴스에 각 http method 에 사용할 Middleware 를 등록하는 메서드입니다.
//methods 는 사용가능한 http method 의 이름과 같아야 합니다.
func (r *Resource) AddMiddleware(middleware gin.HandlerFunc, methods ...string) {
	for _, m := range methods {
		m = strings.ToUpper(m)
		middlewares := r.Middlewares[m]
		middlewares = append(middlewares, middleware)
		r.Middlewares[m] = middlewares
	}
}
