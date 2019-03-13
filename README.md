# gin-restful
gin-restful은 gin으로 restful api를 편하게 개발하기 위해 만든 라이브러리입니다.  
Api에 Resource를 등록하는 형태로 restful api 를 개발할 수 있습니다.  
문서나 주석은 영어를 잘 하지 못해서 한글로 작성했습니다.

## 개요
gin을 이용한 더 편한 restful api 개발을 위해 
Api에 Resource를 등록시키는 형태로 개발할 수 있게 만들었습니다.  
Resource는 어떤 구조체도 될 수 있으며 url이 호출되었을 때 
http method와 이름이 같은 메서드가 호출 됩니다.  
Resource를 등록만 하면 자동으로 Handler Method의 인자를 
분석하여 url를 만들어 gin에 등록합니다.  

## 예시
```go
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
	SampleResource 의 Url 로 GET 요청이 들어왔을 때 실행되는 Handler 입니다.
	gin.H 과 status code 를 반환합니다.
    요청으로 받은 payload를 그대로 반환합니다.
 */
func (r SampleResource) Post(c *gin.Context) (gin.H, int) {
	json := make(map[string]string)
	if err := c.ShouldBindJSON(&json); err != nil {
		return gin.H{"message": err.Error()}, 400
	}
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

```