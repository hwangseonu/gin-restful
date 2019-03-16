package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
	"net/http"
)

/*
	restful api 를 구현하기 위한 SampleResource 구조체입니다.
	gin-restful 의 Resource 구조체 포인터를 임베딩합니다.
	각 Handler 마다 다른 Middleware 들을 적용할 수 있습니다.
 */
type SampleResource struct {
	*gin_restful.Resource
}

type Data struct {
	Name string `json:"name"`
}

/*
	SampleResource 의 Url 로 GET 요청이 들어왔을 때 실행되는 Handler 입니다.
	gin.H 와 status code 를 반환합니다.
	path variable 인 name 을 json 에 담아 반환합니다.
 */
func (r SampleResource) Get(name string) (gin.H, int) {
	return gin.H{
		"name": name,
	}, http.StatusOK
}

/*
	SampleResource 의 Url 로 POST 요청이 들어왔을 때 실행되는 Handler 입니다.
	gin.H 과 status code 를 반환합니다.
    요청으로 받은 payload를 그대로 반환합니다.
 */
func (r SampleResource) Post(c *gin.Context, json Data) (Data, int) {
	return json, 200
}

/*
	Middleware 테스트용 Sample Middleware 입니다.
	콘솔에 "Hello, World" 를 출력합니다.
 */
func SampleMiddleware(c *gin.Context) {
	println("Hello, World")
}

/*
	Api 인스턴스를 "/" 주소로 생성합니다.
	SampleResource 의 인스턴스를 생성하여 Api 에 등록합니다.
	SampleResource GET handler 에 SampleMiddleware 를 등록합니다.
	gin 서버를 5000 포트에서 실행합니다.
 */
func main() {
	r := gin.Default()
	v1 := gin_restful.NewApi(r, "/")
	res := SampleResource{gin_restful.InitResource()}
	res.AddMiddleware(SampleMiddleware, http.MethodGet)
	v1.AddResource(res, "/samples")
	_ = r.Run(":5000")
}
