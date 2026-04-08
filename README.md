# gin-restful

[![Go Reference](https://pkg.go.dev/badge/github.com/hwangseonu/gin-restful.svg)](https://pkg.go.dev/github.com/hwangseonu/gin-restful)
[![CodeFactor](https://www.codefactor.io/repository/github/hwangseonu/gin-restful/badge)](https://www.codefactor.io/repository/github/hwangseonu/gin-restful)

A Flask-RESTful inspired library for building RESTful APIs with [Gin](https://github.com/gin-gonic/gin). Implement only the HTTP methods you need — routes are registered automatically.

## Features

* **Interface Segregation** — Six independent interfaces (`Lister`, `Getter`, `Poster`, `Putter`, `Patcher`, `Deleter`). Only implemented methods become routes.
* **Type-Safe Binding** — Generic `Bind[T]` and `MustBind[T]` helpers for compile-time safe request body parsing.
* **Structured Error Handling** — `Abort()` helper and `HTTPError` type for clean error responses, inspired by Flask-RESTful's `abort()`.
* **Gin Compatible** — Works with `gin.Engine`, `gin.RouterGroup`, middleware, and all standard Gin features.

## Installation

```bash
go get github.com/hwangseonu/gin-restful
```

## Quick Start

```go
package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/hwangseonu/gin-restful"
)

type CreateTodoReq struct {
    Title string `json:"title" binding:"required"`
}

type Todo struct {
    ID    string `json:"id"`
    Title string `json:"title"`
}

type TodoResource struct {
    todos map[string]Todo
}

// GET /todos — implement Lister
func (r *TodoResource) List(c *gin.Context) (any, int, error) {
    items := make([]Todo, 0)
    for _, t := range r.todos {
        items = append(items, t)
    }
    return items, http.StatusOK, nil
}

// GET /todos/:id — implement Getter
func (r *TodoResource) Get(id string, c *gin.Context) (any, int, error) {
    todo, ok := r.todos[id]
    if !ok {
        return nil, 0, restful.Abort(http.StatusNotFound, "todo not found")
    }
    return todo, http.StatusOK, nil
}

// POST /todos — implement Poster
func (r *TodoResource) Post(c *gin.Context) (any, int, error) {
    body, err := restful.Bind[CreateTodoReq](c)
    if err != nil {
        return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
    }
    todo := Todo{ID: "1", Title: body.Title}
    r.todos[todo.ID] = todo
    return todo, http.StatusCreated, nil
}

// DELETE /todos/:id — implement Deleter
func (r *TodoResource) Delete(id string, c *gin.Context) (any, int, error) {
    delete(r.todos, id)
    return nil, http.StatusNoContent, nil
}

// PUT and PATCH are not implemented — those routes won't be registered.

func main() {
    engine := gin.Default()
    api := restful.NewAPI(engine, "/api/v1")

    api.AddResource("/todos", &TodoResource{todos: make(map[string]Todo)})

    log.Fatal(engine.Run(":8080"))
}
```

### Generated Routes

`AddResource("/todos", resource)` registers only the routes for implemented interfaces:

| Method   | Endpoint             | Interface  |
|----------|----------------------|------------|
| `GET`    | `/api/v1/todos`      | `Lister`   |
| `GET`    | `/api/v1/todos/:id`  | `Getter`   |
| `POST`   | `/api/v1/todos`      | `Poster`   |
| `DELETE` | `/api/v1/todos/:id`  | `Deleter`  |

`PUT` and `PATCH` are not registered because the resource doesn't implement `Putter` or `Patcher`.

## Interfaces

Each HTTP method has its own interface. Implement only what you need:

```go
type Lister  interface { List(c *gin.Context) (any, int, error) }              // GET    /path
type Getter  interface { Get(id string, c *gin.Context) (any, int, error) }    // GET    /path/:id
type Poster  interface { Post(c *gin.Context) (any, int, error) }              // POST   /path
type Putter  interface { Put(id string, c *gin.Context) (any, int, error) }    // PUT    /path/:id
type Patcher interface { Patch(id string, c *gin.Context) (any, int, error) }  // PATCH  /path/:id
type Deleter interface { Delete(id string, c *gin.Context) (any, int, error) } // DELETE /path/:id
```

All methods return `(any, int, error)` — the response body, HTTP status code, and an optional error.

## Request Binding

### Bind

Returns the parsed body and an error for custom handling:

```go
body, err := restful.Bind[CreateTodoReq](c)
if err != nil {
    return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
}
```

### MustBind

Automatically responds with `400 Bad Request` and aborts on failure:

```go
body := restful.MustBind[CreateTodoReq](c)
if body == nil {
    return nil, 0, nil // already aborted
}
```

Both use Gin's `ShouldBind`, so the binding method is determined by the `Content-Type` header (JSON, form, XML, etc.).

## Error Handling

Use `Abort` to return structured error responses:

```go
return nil, 0, restful.Abort(http.StatusNotFound, "resource not found")
// Response: 404 {"message": "resource not found"}
```

- `*HTTPError` errors produce their status code and message as JSON.
- Other errors produce `500 {"message": "internal server error"}` — internal details are never leaked to clients.

## Working with Gin Middleware

`NewAPI` accepts `gin.IRouter`, so it works with route groups and middleware:

```go
engine := gin.Default()

// Public routes — regular Gin handlers
engine.GET("/health", healthHandler)

// Protected routes — restful resources with auth middleware
authorized := engine.Group("/", authMiddleware())
api := restful.NewAPI(authorized, "/api/v1")
api.AddResource("/tasks", &TaskResource{})
```

See [`example/hybrid/`](example/hybrid/) for a complete example.

## Examples

| Example | Description |
|---------|-------------|
| [`example/basic/`](example/basic/) | Simple CRUD with List, Get, Post, Put, Delete |
| [`example/complex/`](example/complex/) | Multiple resources, search, pagination, PUT/PATCH separation, relation validation |
| [`example/hybrid/`](example/hybrid/) | RESTful resources + regular Gin handlers, auth middleware, route groups |

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
