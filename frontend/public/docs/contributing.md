# Contributing

Go Initializer is open source and contributions are welcome. This page explains how to add features, fix bugs, and submit pull requests.

## Repository layout

```
.
├── cmd/
│   ├── goini/          # goini CLI binary
│   └── server/         # HTTP server binary
├── frontend/           # React + TypeScript UI
├── internal/
│   ├── generator/      # Core scaffolding engine
│   └── server/         # HTTP handlers and router
├── Makefile            # Developer shortcuts
└── docker-compose.yml
```

## Development setup

```bash
git clone https://github.com/neo7337/go-initializer.git
cd go-initializer

# Install Go dependencies
go mod tidy

# Start the backend in watch mode (requires Air or just use go run)
go run ./cmd/server/...

# In a second terminal, start the frontend
cd frontend && npm install && npm start
```

---

## Adding a new project type

1. Create a file `internal/generator/gen_<type>.go` and implement the `Generator` interface:

```go
type Generator interface {
    Generate(request CreateProjectRequest) (*bytes.Buffer, error)
}
```

2. Register it in `internal/generator/registry.go`:

```go
var GeneratorRegistry = map[string]Generator{
    // existing entries...
    "my-new-type": &MyNewTypeGenerator{},
}
```

3. Add it to the supported maps in the same file:

```go
var SupportedProjectTypesLabelsMap = map[string]string{
    // existing entries...
    "my-new-type": "My New Type",
}
```

4. Add tests in `internal/generator/generators_test.go`.

---

## Adding a new framework

1. Add the framework to the `SupportedFrameworksMap` entry for the relevant project type in `registry.go`.
2. Add its module dependency to `DependencyMap`.
3. Update the corresponding generator to handle the new framework value.

---

## Adding a new addon

1. Create `internal/generator/addon_<category>.go` implementing `AddonGenerator`:

```go
type AddonGenerator interface {
    Generate(folderName string, addons []string, zw *zip.Writer) error
}
```

2. Register it in `addonRegistry` in `registry.go`.
3. Add the valid values to `SupportedAddonsMap`.

---

## Running tests

```bash
# All Go tests
go test ./...

# With race detector
go test -race ./...

# Frontend type-check
cd frontend && npx tsc --noEmit
```

---

## Code style

- Go code follows standard `gofmt` formatting. Run `go fmt ./...` before committing.
- TypeScript follows the existing code style — no extra linter config is required.
- Keep PRs focused: one feature or bug fix per PR makes review easier.

---

## Submitting a pull request

1. Fork the repository and create a feature branch
2. Make your changes and add tests
3. Run `go test ./...` to ensure nothing is broken
4. Push your branch and open a PR against `main`
5. Fill in the PR template and wait for a review

---

## Reporting bugs

Open an issue at [github.com/neo7337/go-initializer/issues](https://github.com/neo7337/go-initializer/issues) with:

- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected vs actual behaviour
