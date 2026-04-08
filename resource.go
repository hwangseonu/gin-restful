package gin_restful

import "github.com/gin-gonic/gin"

// 각 HTTP 메서드별 독립 인터페이스. 필요한 것만 구현하면 해당 라우트만 등록.

type Lister interface {
	List(c *gin.Context) (any, int, error)
}

type Getter interface {
	Get(id string, c *gin.Context) (any, int, error)
}

type Poster interface {
	Post(c *gin.Context) (any, int, error)
}

type Putter interface {
	Put(id string, c *gin.Context) (any, int, error)
}

type Patcher interface {
	Patch(id string, c *gin.Context) (any, int, error)
}

type Deleter interface {
	Delete(id string, c *gin.Context) (any, int, error)
}
