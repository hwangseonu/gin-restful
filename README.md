# gin-restful

A Go library that simplifies and accelerates RESTful API development using the Gin framework. It abstracts away repetitive routing and handler setups, allowing you to easily implement **Create, Read, Update, Delete (CRUD)** functionalities.

---

## Key Features

* **Automatic Routing:** Automatically maps HTTP methods (POST, GET, PUT, DELETE) to the corresponding CRUD functions defined in the `restful.Resource` interface.
* **Interface-Based Development:** Provides a consistent way to define API resources by simply implementing a predefined interface.
* **Easy Schema Binding:** The `RequestBody` method allows you to easily define the JSON schema for incoming request bodies.

---

## Getting Started

### Installation

First, install the Gin framework and the `gin-restful` library.

```bash
go get github.com/gin-gonic/gin
go get github.com/hwangseonu/gin-restful
```

### Usage Example
The following code is a complete example of how to use the gin-restful library to create a simple sample API.

```go
package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	restful "github.com/hwangseonu/gin-restful"
)

// Define the schema for the request body.
type SampleSchema struct {
	Message string `json:"message"`
}

var offset = 0

// Define the struct that will implement the restful.Resource interface.
type Sample struct {
	database map[string]SampleSchema
}

// Returns the schema for the request body.
func (r *Sample) RequestBody(_ string) any {
	return new(SampleSchema)
}

// Creates a new resource (POST).
func (r *Sample) Create(body interface{}, _ *gin.Context) (gin.H, int, error) {
	sample := body.(*SampleSchema)
	id := strconv.Itoa(offset)
	offset += 1
	r.database[id] = *sample

	return gin.H{
		"message": sample.Message,
	}, http.StatusCreated, nil
}

// Reads a specific resource (GET).
func (r *Sample) Read(id string, _ *gin.Context) (gin.H, int, error) {
	sample, ok := r.database[id]

	if !ok {
		return gin.H{}, 404, nil
	} else {
		return gin.H{"message": sample.Message}, http.StatusOK, nil
	}
}

// Reads all resources (GET).
func (r *Sample) ReadAll(_ *gin.Context) (gin.H, int, error) {
	samples := make(map[string]gin.H)
	for k, v := range r.database {
		samples[k] = gin.H{"message": v.Message}
	}

	return gin.H{
		"samples": samples,
	}, http.StatusOK, nil
}

// Updates a specific resource (PUT).
func (r *Sample) Update(id string, body interface{}, _ *gin.Context) (gin.H, int, error) {
	sample := body.(*SampleSchema)

	r.database[id] = *sample

	return gin.H{}, http.StatusNoContent, nil
}

// Deletes a specific resource (DELETE).
func (r *Sample) Delete(id string, _ *gin.Context) (gin.H, int, error) {
	delete(r.database, id)
	return gin.H{}, http.StatusNoContent, nil
}

func main() {
	// Initialize the Gin engine and restful API.
	engine := gin.Default()
	api := restful.NewAPI("/api/v1")

	// Create a resource instance.
	sample := &Sample{make(map[string]SampleSchema)}

	// Register the resource routing.
	api.RegisterResource("/samples", sample)
	api.RegisterHandlers(engine)

	// Run the server.
	e := engine.Run(":8080")
	if e != nil {
		log.Fatalln(e)
	}
}
```

### API Endpoints
The api.RegisterResource("/samples", sample) line in the example code automatically generates the following RESTful endpoints.

| HTTP Method | Endpoint            | Description                  |
|:------------| :------------------ | :--------------------------- |
| `POST`      | `/api/v1/samples`   | Creates a new resource       |
| `GET`       | `/api/v1/samples`   | Retrieves all resources      |
| `GET`       | `/api/v1/samples/:id` | Retrieves a specific resource |
| `PUT`       | `/api/v1/samples/:id` | Updates a specific resource   |
| `PATCH`     | `/api/v1/samples/:id` | Updates a specific resource   |
| `DELETE`    | `/api/v1/samples/:id` | Deletes a specific resource   |

## License
This project is licensed under the MIT License. See the LICENSE file for details.