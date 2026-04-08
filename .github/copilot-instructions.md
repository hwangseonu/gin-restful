# Copilot Instructions for gin-restful

Flask-RESTful style REST API library for Gin. See AGENTS.md for full details.

## Resource Pattern

Implement independent interfaces for each HTTP method — only implemented methods become routes:

```go
// Implement Lister for GET /path
func (r *MyResource) List(c *gin.Context) (any, int, error) { ... }

// Implement Getter for GET /path/:id
func (r *MyResource) Get(id string, c *gin.Context) (any, int, error) { ... }

// Implement Poster for POST /path
func (r *MyResource) Post(c *gin.Context) (any, int, error) { ... }

// Implement Putter for PUT /path/:id
func (r *MyResource) Put(id string, c *gin.Context) (any, int, error) { ... }

// Implement Patcher for PATCH /path/:id
func (r *MyResource) Patch(id string, c *gin.Context) (any, int, error) { ... }

// Implement Deleter for DELETE /path/:id
func (r *MyResource) Delete(id string, c *gin.Context) (any, int, error) { ... }
```

## Request Binding

```go
// Option 1: Bind — custom error handling
body, err := restful.Bind[CreateReq](c)
if err != nil {
    return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
}

// Option 2: MustBind — automatic 400 abort
body := restful.MustBind[CreateReq](c)
if body == nil {
    return nil, 0, nil
}
```

## Error Handling

```go
// Use Abort for client errors
return nil, 0, restful.Abort(404, "not found")

// Never use fmt.Errorf for client-facing errors — it produces 500 with generic message
```

## Registration

```go
api := restful.NewAPI(engine, "/api/v1")
api.AddResource("/items", &ItemResource{})
```

## Conventions

- All handler methods return `(any, int, error)`
- Use `sync.RWMutex` for concurrent map/state access in resource structs
- Package name is `restful` — no import alias needed
- Run `go test ./... -v` to verify changes
