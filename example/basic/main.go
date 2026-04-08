package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
)

type SampleSchema struct {
	Message string `json:"message" binding:"required"`
}

type Sample struct {
	mu       sync.RWMutex
	database map[string]SampleSchema
	nextID   int
}

func (r *Sample) List(c *gin.Context) (any, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	samples := make(map[string]gin.H)
	for k, v := range r.database {
		samples[k] = gin.H{"message": v.Message}
	}
	return gin.H{"samples": samples}, http.StatusOK, nil
}

func (r *Sample) Get(id string, c *gin.Context) (any, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

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

	r.mu.Lock()
	defer r.mu.Unlock()

	id := strconv.Itoa(r.nextID)
	r.nextID++
	r.database[id] = *body
	return gin.H{"message": body.Message}, http.StatusCreated, nil
}

func (r *Sample) Put(id string, c *gin.Context) (any, int, error) {
	body, err := restful.Bind[SampleSchema](c)
	if err != nil {
		return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.database[id] = *body
	return nil, http.StatusNoContent, nil
}

func (r *Sample) Delete(id string, c *gin.Context) (any, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

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
