package main

import (
	"log"

	"github.com/gin-gonic/gin"
	restful "github.com/hwangseonu/gin-restful"
)

type SampleSchema struct {
	Message string `json:"message"`
}

type Sample struct {
	restful.Resource
}

func (r Sample) RequestBody() interface{} {
	return new(SampleSchema)
}

func (r Sample) Create(body interface{}, _ *gin.Context) (gin.H, error) {
	sample := body.(*SampleSchema)

	return gin.H{
		"message": sample.Message,
	}, nil
}

func main() {
	engine := gin.Default()
	api := restful.NewAPI("/api/v1")

	sample := Sample{}

	api.RegisterResource("/samples", sample)
	api.RegisterHandlers(engine)

	e := engine.Run(":8080")
	if e != nil {
		log.Fatalln(e)
	}
}
