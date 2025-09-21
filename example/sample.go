package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	restful "github.com/hwangseonu/gin-restful"
)

type SampleSchema struct {
	Message string `json:"message"`
}

var offset = 0

type Sample struct {
	database map[string]SampleSchema
}

func (r *Sample) RequestBody(method string) any {
	if method == "GET" {
		return nil
	}
	return new(SampleSchema)
}

func (r *Sample) Create(body interface{}, _ *gin.Context) (gin.H, int, error) {
	sample := body.(*SampleSchema)
	id := strconv.Itoa(offset)
	offset += 1
	r.database[id] = *sample

	return gin.H{
		"message": sample.Message,
	}, http.StatusCreated, nil
}

func (r *Sample) Read(id string, _ *gin.Context) (gin.H, int, error) {
	sample, ok := r.database[id]

	if !ok {
		return gin.H{}, 404, nil
	}
	return gin.H{"message": sample.Message}, http.StatusOK, nil
}

func (r *Sample) ReadAll(_ *gin.Context) (gin.H, int, error) {
	samples := make(map[string]gin.H)
	for k, v := range r.database {
		samples[k] = gin.H{"message": v.Message}
	}

	return gin.H{
		"samples": samples,
	}, http.StatusOK, nil
}

func (r *Sample) Update(id string, body interface{}, _ *gin.Context) (gin.H, int, error) {
	sample := body.(*SampleSchema)

	r.database[id] = *sample

	return gin.H{}, http.StatusNoContent, nil
}

func (r *Sample) Delete(id string, _ *gin.Context) (gin.H, int, error) {
	delete(r.database, id)
	return gin.H{}, http.StatusNoContent, nil
}

func main() {
	engine := gin.Default()
	api := restful.NewAPI("/api/v1")

	sample := &Sample{make(map[string]SampleSchema)}

	api.RegisterResource("/samples", sample)
	api.RegisterHandlers(&engine.RouterGroup)

	e := engine.Run(":8080")
	if e != nil {
		log.Fatalln(e)
	}
}
