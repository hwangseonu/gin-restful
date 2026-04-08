package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	restful "github.com/hwangseonu/gin-restful"
)

type SampleSchema struct {
	Message string `json:"message" binding:"required"`
}

var offset = 0

type Sample struct {
	database map[string]SampleSchema
}

func (r *Sample) List(c *gin.Context) (any, int, error) {
	samples := make(map[string]gin.H)
	for k, v := range r.database {
		samples[k] = gin.H{"message": v.Message}
	}
	return gin.H{"samples": samples}, http.StatusOK, nil
}

func (r *Sample) Get(id string, c *gin.Context) (any, int, error) {
	sample, ok := r.database[id]
	if !ok {
		return nil, 0, restful.Abort(http.StatusNotFound, "sample not found")
	}
	return gin.H{"message": sample.Message}, http.StatusOK, nil
}

func (r *Sample) Post(c *gin.Context) (any, int, error) {
	body, err := restful.Bind[SampleSchema](c)
	if err != nil {
		return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
	}
	id := strconv.Itoa(offset)
	offset++
	r.database[id] = *body
	return gin.H{"message": body.Message}, http.StatusCreated, nil
}

func (r *Sample) Put(id string, c *gin.Context) (any, int, error) {
	body, err := restful.Bind[SampleSchema](c)
	if err != nil {
		return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
	}
	r.database[id] = *body
	return nil, http.StatusNoContent, nil
}

func (r *Sample) Delete(id string, c *gin.Context) (any, int, error) {
	delete(r.database, id)
	return nil, http.StatusNoContent, nil
}

func main() {
	engine := gin.Default()
	api := restful.NewAPI(engine, "/api/v1")

	sample := &Sample{database: make(map[string]SampleSchema)}

	api.AddResource("/samples", sample)

	if err := engine.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}
