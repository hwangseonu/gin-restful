# AGENTS.md

This file provides guidance to AI coding agents working with this repository.

## Project Overview

`gin-restful` is a Go library that provides a Flask-RESTful style abstraction over the Gin web framework. Resources are defined by implementing independent HTTP method interfaces — only implemented methods become routes.

## Commands

```bash
go build ./...              # Build the library
go vet ./...                # Lint
go test ./... -v            # Run tests
go build -o /tmp/bin ./example/basic/  # Build example (use -o to avoid dir conflict)
```

## Architecture

### Core Files

| File | Purpose |
|------|---------|
| `resource.go` | 6 interfaces: `Lister`, `Getter`, `Poster`, `Putter`, `Patcher`, `Deleter` |
| `api.go` | `API` struct, `NewAPI`, `AddResource` — type-asserts resource and registers Gin routes |
| `bind.go` | `Bind[T]` and `MustBind[T]` — generic request body binding helpers |
| `errors.go` | `HTTPError` type and `Abort()` helper for structured error responses |
| `handler.go` | `makeHandler` (response dispatcher), `handleError`, `normalizePath` |

### Design Patterns

- **Interface segregation**: Each HTTP method is a separate interface. Implement only what you need.
- **Type-safe binding**: `Bind[T]` / `MustBind[T]` use Go generics for compile-time type checking on request bodies. Response type is `any`.
- **Error handling**: Return `*HTTPError` (via `Abort()`) for client errors with status+message. Non-HTTPError errors produce `500` with a generic message — internal details are never leaked.
- **204 No Content**: When status is 204, the framework sends no response body.
- **Abort on error**: All error paths call `c.AbortWithStatusJSON` to stop the middleware chain.

### Package

Package name is `restful`. No import alias needed.

## Code Conventions

- Handler methods return `(any, int, error)` — response body, HTTP status code, error.
- Use `Bind[T]` when you need custom error handling, `MustBind[T]` for automatic 400 abort.
- Use `Abort(status, message)` to return HTTP errors — never return raw `fmt.Errorf` for client-facing errors.
- `AddResource` panics if the resource implements none of the 6 interfaces — this catches typos/signature mismatches early.
- All examples use `sync.RWMutex` for concurrent map access — follow this pattern.

## Examples

| Example | Path | Description |
|---------|------|-------------|
| Basic | `example/basic/` | Simple CRUD |
| Complex | `example/complex/` | Multiple resources, pagination, PUT/PATCH split |
| Hybrid | `example/hybrid/` | RESTful + regular Gin handlers, auth middleware |
