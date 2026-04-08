# Contributing to gin-restful

Thank you for your interest in contributing!

## Development Setup

```bash
git clone https://github.com/hwangseonu/gin-restful.git
cd gin-restful
go test ./... -v
```

## Making Changes

1. Create a feature branch from `master`.
2. Write tests first (TDD) — tests should fail before implementation.
3. Implement the minimal code to pass tests.
4. Run the full test suite with race detection:
   ```bash
   go test ./... -v -race
   go vet ./...
   ```
5. Update examples if the public API changes.

## Pull Request Guidelines

- Keep PRs focused — one feature or fix per PR.
- Include tests for new functionality.
- Update documentation (README, GoDoc, AGENTS.md) if the public API changes.
- Commit messages should describe *why*, not just *what*.

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`).
- All public types and functions must have GoDoc comments.
- Use `Abort()` for client-facing errors, never raw `fmt.Errorf`.
- Use `sync.RWMutex` for shared state in examples and resource structs.
- Handler methods return `(any, int, error)`.

## Running Benchmarks

```bash
go test -bench=. -benchmem ./...
```
