// Package restful provides a Flask-RESTful inspired abstraction over the
// Gin web framework. Resources are defined by implementing one or more of the
// HTTP method interfaces (Lister, Getter, Poster, Putter, Patcher, Deleter).
// Only implemented interfaces are registered as routes.
package restful

import "github.com/gin-gonic/gin"

// Lister handles GET requests on a collection path (e.g. GET /items).
type Lister interface {
	List(c *gin.Context) (any, int, error)
}

// Getter handles GET requests for a single resource (e.g. GET /items/:id).
type Getter interface {
	Get(id string, c *gin.Context) (any, int, error)
}

// Poster handles POST requests to create a new resource (e.g. POST /items).
type Poster interface {
	Post(c *gin.Context) (any, int, error)
}

// Putter handles PUT requests to replace a resource (e.g. PUT /items/:id).
type Putter interface {
	Put(id string, c *gin.Context) (any, int, error)
}

// Patcher handles PATCH requests to partially update a resource (e.g. PATCH /items/:id).
type Patcher interface {
	Patch(id string, c *gin.Context) (any, int, error)
}

// Deleter handles DELETE requests to remove a resource (e.g. DELETE /items/:id).
type Deleter interface {
	Delete(id string, c *gin.Context) (any, int, error)
}
