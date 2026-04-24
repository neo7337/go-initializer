package generator

import (
	"archive/zip"
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// zipEntries reads all entries from the zip buffer and returns a map from entry
// path to file contents.
func zipEntries(t *testing.T, buf *bytes.Buffer) map[string]string {
	t.Helper()
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	require.NoError(t, err, "zip.NewReader failed")
	entries := make(map[string]string, len(r.File))
	for _, f := range r.File {
		rc, err := f.Open()
		require.NoError(t, err)
		var content bytes.Buffer
		_, err = content.ReadFrom(rc)
		rc.Close()
		require.NoError(t, err)
		entries[f.Name] = content.String()
	}
	return entries
}

// hasEntry asserts that the zip contains the entry at key and returns its content.
func hasEntry(t *testing.T, entries map[string]string, key string) string {
	t.Helper()
	v, ok := entries[key]
	assert.True(t, ok, "expected zip entry %q to exist; available: %v", key, zipKeys(entries))
	return v
}

func zipKeys(entries map[string]string) []string {
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	return keys
}

// ─── SimpleProjectGenerator ──────────────────────────────────────────────────

func TestSimpleProjectGenerator_Basic(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "simple-project",
		GoVersion:   "1.24.6",
		Framework:   "golly",
		ModuleName:  "github.com/acme/myproj",
		Name:        "myproj",
	}
	buf, err := GeneratorRegistry["simple-project"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)

	hasEntry(t, entries, "myproj/README.md")
	hasEntry(t, entries, "myproj/go.mod")
	hasEntry(t, entries, "myproj/cmd/myproj/main.go")
	hasEntry(t, entries, "myproj/internal/service.go")
	hasEntry(t, entries, "myproj/Makefile")
	hasEntry(t, entries, "myproj/.gitignore")
}

func TestSimpleProjectGenerator_ReadmeContainsName(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "simple-project",
		GoVersion:   "1.24.6",
		Framework:   "golly",
		ModuleName:  "github.com/acme/alpha",
		Name:        "alpha",
		Description: "my alpha service",
	}
	buf, err := GeneratorRegistry["simple-project"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	readme := hasEntry(t, entries, "alpha/README.md")
	assert.Contains(t, readme, "alpha")
	assert.Contains(t, readme, "my alpha service")
}

func TestSimpleProjectGenerator_WithDocker(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType:   "simple-project",
		GoVersion:     "1.24.6",
		Framework:     "golly",
		ModuleName:    "github.com/acme/myproj",
		Name:          "myproj",
		DockerSupport: true,
	}
	buf, err := GeneratorRegistry["simple-project"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "myproj/Dockerfile")
}

func TestSimpleProjectGenerator_WithoutDocker(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType:   "simple-project",
		GoVersion:     "1.24.6",
		Framework:     "golly",
		ModuleName:    "github.com/acme/myproj",
		Name:          "myproj",
		DockerSupport: false,
	}
	buf, err := GeneratorRegistry["simple-project"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	_, ok := entries["myproj/Dockerfile"]
	assert.False(t, ok, "Dockerfile should not be included when DockerSupport=false")
}

func TestSimpleProjectGenerator_WithRedisAddon(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "simple-project",
		GoVersion:   "1.24.6",
		Framework:   "golly",
		ModuleName:  "github.com/acme/myproj",
		Name:        "myproj",
		Addons:      map[string][]string{"cache": {"redis"}},
	}
	buf, err := GeneratorRegistry["simple-project"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	content := hasEntry(t, entries, "myproj/internal/cache/cache.go")
	assert.Contains(t, content, "NewRedisClient")
}

func TestSimpleProjectGenerator_WithGormAddon(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "simple-project",
		GoVersion:   "1.24.6",
		Framework:   "golly",
		ModuleName:  "github.com/acme/myproj",
		Name:        "myproj",
		Addons:      map[string][]string{"database": {"gorm"}},
	}
	buf, err := GeneratorRegistry["simple-project"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	content := hasEntry(t, entries, "myproj/internal/database/database.go")
	assert.Contains(t, content, "NewGormDB")
}

func TestSimpleProjectGenerator_WithZapLogging(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "simple-project",
		GoVersion:   "1.24.6",
		Framework:   "golly",
		ModuleName:  "github.com/acme/myproj",
		Name:        "myproj",
		Addons:      map[string][]string{"other": {"zap"}},
	}
	buf, err := GeneratorRegistry["simple-project"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	content := hasEntry(t, entries, "myproj/internal/logger/logger.go")
	assert.Contains(t, content, "zap")
}

func TestSimpleProjectGenerator_EmptyNameDefaultsToMyproject(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "simple-project",
		GoVersion:   "1.24.6",
		Framework:   "golly",
		ModuleName:  "github.com/acme/myproject",
		Name:        "",
	}
	buf, err := GeneratorRegistry["simple-project"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "myproject/README.md")
}

func TestSimpleProjectGenerator_MainPackageIsMain(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "simple-project",
		GoVersion:   "1.24.6",
		Framework:   "golly",
		ModuleName:  "github.com/acme/p",
		Name:        "p",
	}
	buf, err := GeneratorRegistry["simple-project"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	main := hasEntry(t, entries, "p/cmd/p/main.go")
	assert.Contains(t, main, "package main")
}

// ─── MicroserviceGenerator ───────────────────────────────────────────────────

var microserviceFrameworks = []string{"golly", "gin", "echo", "fiber", "gokit"}

func TestMicroserviceGenerator_RequiredFiles(t *testing.T) {
	for _, fw := range microserviceFrameworks {
		t.Run(fw, func(t *testing.T) {
			req := CreateProjectRequest{
				ProjectType: "microservice",
				GoVersion:   "1.24.6",
				Framework:   fw,
				ModuleName:  "github.com/acme/svc",
				Name:        "svc",
			}
			buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
			require.NoError(t, err)
			entries := zipEntries(t, buf)

			hasEntry(t, entries, "svc/README.md")
			hasEntry(t, entries, "svc/go.mod")
			hasEntry(t, entries, "svc/cmd/svc/main.go")
			hasEntry(t, entries, "svc/internal/handler/handler.go")
			hasEntry(t, entries, "svc/internal/router/router.go")
			hasEntry(t, entries, "svc/internal/service/service.go")
			hasEntry(t, entries, "svc/Makefile")
			hasEntry(t, entries, "svc/.gitignore")
		})
	}
}

func TestMicroserviceGenerator_MainPackageIsMain(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	main := hasEntry(t, entries, "svc/cmd/svc/main.go")
	assert.Contains(t, main, "package main")
}

func TestMicroserviceGenerator_WithDocker(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType:   "microservice",
		GoVersion:     "1.24.6",
		Framework:     "gin",
		ModuleName:    "github.com/acme/svc",
		Name:          "svc",
		DockerSupport: true,
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "svc/Dockerfile")
}

func TestMicroserviceGenerator_WithCacheAddon(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
		Addons:      map[string][]string{"cache": {"redis"}},
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "svc/internal/cache/cache.go")
}

func TestMicroserviceGenerator_WithDatabaseAddon(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
		Addons:      map[string][]string{"database": {"ent"}},
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "svc/internal/database/database.go")
}

func TestMicroserviceGenerator_WithLogrusAddon(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
		Addons:      map[string][]string{"other": {"logrus"}},
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	content := hasEntry(t, entries, "svc/internal/logger/logger.go")
	assert.Contains(t, content, "logrus")
}

func TestMicroserviceGenerator_EmptyNameDefaultsToMyservice(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "golly",
		ModuleName:  "github.com/acme/myservice",
		Name:        "",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "myservice/README.md")
}

func TestMicroserviceGenerator_ServiceStubContainsName(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/billing",
		Name:        "billing",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	svc := hasEntry(t, entries, "billing/internal/service/service.go")
	assert.Contains(t, svc, "BillingService")
}

// ─── CLIAppGenerator ─────────────────────────────────────────────────────────

var cliFrameworks = []string{"golly", "cobra", "urfave", "kingpin"}

func TestCLIAppGenerator_RequiredFiles(t *testing.T) {
	for _, fw := range cliFrameworks {
		t.Run(fw, func(t *testing.T) {
			req := CreateProjectRequest{
				ProjectType: "cli-app",
				GoVersion:   "1.24.6",
				Framework:   fw,
				ModuleName:  "github.com/acme/cli",
				Name:        "cli",
			}
			buf, err := GeneratorRegistry["cli-app"].Generate(context.Background(), req)
			require.NoError(t, err)
			entries := zipEntries(t, buf)

			hasEntry(t, entries, "cli/README.md")
			hasEntry(t, entries, "cli/go.mod")
			hasEntry(t, entries, "cli/main.go")
			hasEntry(t, entries, "cli/cmd/root.go")
			hasEntry(t, entries, "cli/cmd/cli.go")
			hasEntry(t, entries, "cli/Makefile")
			hasEntry(t, entries, "cli/.gitignore")
		})
	}
}

func TestCLIAppGenerator_MainDelegatesToCmdExecute(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "cli-app",
		GoVersion:   "1.24.6",
		Framework:   "cobra",
		ModuleName:  "github.com/acme/tool",
		Name:        "tool",
	}
	buf, err := GeneratorRegistry["cli-app"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	main := hasEntry(t, entries, "tool/main.go")
	assert.Contains(t, main, "package main")
	assert.Contains(t, main, "cmd.Execute()")
}

func TestCLIAppGenerator_RootCmdPackage(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "cli-app",
		GoVersion:   "1.24.6",
		Framework:   "cobra",
		ModuleName:  "github.com/acme/tool",
		Name:        "tool",
	}
	buf, err := GeneratorRegistry["cli-app"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	root := hasEntry(t, entries, "tool/cmd/root.go")
	assert.Contains(t, root, "package cmd")
	assert.Contains(t, root, "Execute")
}

func TestCLIAppGenerator_SubCmdPackage(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "cli-app",
		GoVersion:   "1.24.6",
		Framework:   "cobra",
		ModuleName:  "github.com/acme/tool",
		Name:        "tool",
	}
	buf, err := GeneratorRegistry["cli-app"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	sub := hasEntry(t, entries, "tool/cmd/tool.go")
	assert.Contains(t, sub, "package cmd")
}

func TestCLIAppGenerator_MakefileUsesRootPkg(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "cli-app",
		GoVersion:   "1.24.6",
		Framework:   "cobra",
		ModuleName:  "github.com/acme/tool",
		Name:        "tool",
	}
	buf, err := GeneratorRegistry["cli-app"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	mk := hasEntry(t, entries, "tool/Makefile")
	// CLI apps have main package at root "."
	assert.Contains(t, mk, "MAIN_PKG=.")
}

func TestCLIAppGenerator_WithZapLogging(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "cli-app",
		GoVersion:   "1.24.6",
		Framework:   "cobra",
		ModuleName:  "github.com/acme/tool",
		Name:        "tool",
		Addons:      map[string][]string{"other": {"zap"}},
	}
	buf, err := GeneratorRegistry["cli-app"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	content := hasEntry(t, entries, "tool/internal/logger/logger.go")
	assert.Contains(t, content, "zap")
}

func TestCLIAppGenerator_EmptyNameDefaultsToMycli(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "cli-app",
		GoVersion:   "1.24.6",
		Framework:   "cobra",
		ModuleName:  "github.com/acme/mycli",
		Name:        "",
	}
	buf, err := GeneratorRegistry["cli-app"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "mycli/README.md")
}

// ─── Registry coverage ───────────────────────────────────────────────────────

func TestGeneratorRegistry_KnownTypesRegistered(t *testing.T) {
	for _, typ := range []string{"simple-project", "microservice", "cli-app", "api-server", "ai-agent"} {
		_, ok := GeneratorRegistry[typ]
		assert.True(t, ok, "expected %q to be in GeneratorRegistry", typ)
	}
}

// ─── GoMod correctness via generator ─────────────────────────────────────────

func TestMicroserviceGenerator_GoModContainsFrameworkDep(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "echo",
		ModuleName:  "github.com/acme/echosvc",
		Name:        "echosvc",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	gomod := hasEntry(t, entries, "echosvc/go.mod")
	assert.Contains(t, gomod, "github.com/labstack/echo")
}

func TestCLIAppGenerator_GoModContainsCobraDep(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "cli-app",
		GoVersion:   "1.24.6",
		Framework:   "cobra",
		ModuleName:  "github.com/acme/ctool",
		Name:        "ctool",
	}
	buf, err := GeneratorRegistry["cli-app"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	gomod := hasEntry(t, entries, "ctool/go.mod")
	assert.Contains(t, gomod, "github.com/spf13/cobra")
}

// ─── Zip slip guard: folder name must always be valid ────────────────────────

func TestGenerators_AllEntryPathsAreRelative(t *testing.T) {
	generators := map[string]CreateProjectRequest{
		"simple-project": {ProjectType: "simple-project", GoVersion: "1.24.6", Framework: "golly", ModuleName: "m", Name: "safe"},
		"microservice":   {ProjectType: "microservice", GoVersion: "1.24.6", Framework: "gin", ModuleName: "m", Name: "safe"},
		"cli-app":        {ProjectType: "cli-app", GoVersion: "1.24.6", Framework: "cobra", ModuleName: "m", Name: "safe"},
		"api-server":     {ProjectType: "api-server", GoVersion: "1.24.6", Framework: "gin", ModuleName: "m", Name: "safe"},
		"ai-agent":       {ProjectType: "ai-agent", GoVersion: "1.24.6", Framework: "openai", ModuleName: "m", Name: "safe"},
	}
	for typ, req := range generators {
		t.Run(typ, func(t *testing.T) {
			buf, err := GeneratorRegistry[typ].Generate(context.Background(), req)
			require.NoError(t, err)
			entries := zipEntries(t, buf)
			for path := range entries {
				assert.False(t, strings.HasPrefix(path, "/"), "zip entry %q must not be an absolute path", path)
				assert.False(t, strings.HasPrefix(path, ".."), "zip entry %q must not escape root", path)
			}
		})
	}
}

// ─── APIServerGenerator ──────────────────────────────────────────────────────

var apiServerFrameworks = []string{"golly", "gin", "echo", "fiber", "chi"}

func TestAPIServerGenerator_RequiredFiles(t *testing.T) {
	for _, fw := range apiServerFrameworks {
		t.Run(fw, func(t *testing.T) {
			req := CreateProjectRequest{
				ProjectType: "api-server",
				GoVersion:   "1.24.6",
				Framework:   fw,
				ModuleName:  "github.com/acme/api",
				Name:        "api",
			}
			buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
			require.NoError(t, err)
			entries := zipEntries(t, buf)

			hasEntry(t, entries, "api/README.md")
			hasEntry(t, entries, "api/go.mod")
			hasEntry(t, entries, "api/cmd/api/main.go")
			hasEntry(t, entries, "api/internal/handler/handler.go")
			hasEntry(t, entries, "api/internal/router/router.go")
			hasEntry(t, entries, "api/internal/service/service.go")
			hasEntry(t, entries, "api/Makefile")
			hasEntry(t, entries, "api/.gitignore")
		})
	}
}

func TestAPIServerGenerator_MainPackageIsMain(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "api-server",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/api",
		Name:        "api",
	}
	buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	main := hasEntry(t, entries, "api/cmd/api/main.go")
	assert.Contains(t, main, "package main")
}

func TestAPIServerGenerator_WithDocker(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType:   "api-server",
		GoVersion:     "1.24.6",
		Framework:     "chi",
		ModuleName:    "github.com/acme/api",
		Name:          "api",
		DockerSupport: true,
	}
	buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "api/Dockerfile")
}

func TestAPIServerGenerator_WithoutDocker(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType:   "api-server",
		GoVersion:     "1.24.6",
		Framework:     "chi",
		ModuleName:    "github.com/acme/api",
		Name:          "api",
		DockerSupport: false,
	}
	buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	_, ok := entries["api/Dockerfile"]
	assert.False(t, ok, "Dockerfile should not be included when DockerSupport=false")
}

func TestAPIServerGenerator_GoModContainsFrameworkDep(t *testing.T) {
	cases := map[string]string{
		"gin":   "github.com/gin-gonic/gin",
		"echo":  "github.com/labstack/echo",
		"fiber": "github.com/gofiber/fiber",
		"chi":   "github.com/go-chi/chi",
	}
	for fw, dep := range cases {
		t.Run(fw, func(t *testing.T) {
			req := CreateProjectRequest{
				ProjectType: "api-server",
				GoVersion:   "1.24.6",
				Framework:   fw,
				ModuleName:  "github.com/acme/api",
				Name:        "api",
			}
			buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
			require.NoError(t, err)
			entries := zipEntries(t, buf)
			gomod := hasEntry(t, entries, "api/go.mod")
			assert.Contains(t, gomod, dep)
		})
	}
}

func TestAPIServerGenerator_WithCacheAddon(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "api-server",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/api",
		Name:        "api",
		Addons:      map[string][]string{"cache": {"redis"}},
	}
	buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	content := hasEntry(t, entries, "api/internal/cache/cache.go")
	assert.Contains(t, content, "NewRedisClient")
}

func TestAPIServerGenerator_WithDatabaseAddon(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "api-server",
		GoVersion:   "1.24.6",
		Framework:   "echo",
		ModuleName:  "github.com/acme/api",
		Name:        "api",
		Addons:      map[string][]string{"database": {"gorm"}},
	}
	buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	content := hasEntry(t, entries, "api/internal/database/database.go")
	assert.Contains(t, content, "NewGormDB")
}

func TestAPIServerGenerator_WithLoggingAddon(t *testing.T) {
	for _, addon := range []string{"zap", "logrus"} {
		t.Run(addon, func(t *testing.T) {
			req := CreateProjectRequest{
				ProjectType: "api-server",
				GoVersion:   "1.24.6",
				Framework:   "fiber",
				ModuleName:  "github.com/acme/api",
				Name:        "api",
				Addons:      map[string][]string{"other": {addon}},
			}
			buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
			require.NoError(t, err)
			entries := zipEntries(t, buf)
			content := hasEntry(t, entries, "api/internal/logger/logger.go")
			assert.Contains(t, content, addon)
		})
	}
}

func TestAPIServerGenerator_ServiceStubContainsName(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "api-server",
		GoVersion:   "1.24.6",
		Framework:   "chi",
		ModuleName:  "github.com/acme/orders",
		Name:        "orders",
	}
	buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	svc := hasEntry(t, entries, "orders/internal/service/service.go")
	assert.Contains(t, svc, "OrdersService")
}

func TestAPIServerGenerator_MakefileUsesCmdPkg(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "api-server",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/api",
		Name:        "api",
	}
	buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	mk := hasEntry(t, entries, "api/Makefile")
	assert.Contains(t, mk, "MAIN_PKG=./cmd/api")
}

func TestAPIServerGenerator_EmptyNameDefaultsToMyapiserver(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "api-server",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/myapiserver",
		Name:        "",
	}
	buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "myapiserver/README.md")
}

func TestAPIServerGenerator_ReadmeContainsNameAndDescription(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "api-server",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/myapi",
		Name:        "myapi",
		Description: "a REST API server",
	}
	buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	readme := hasEntry(t, entries, "myapi/README.md")
	assert.Contains(t, readme, "myapi")
	assert.Contains(t, readme, "a REST API server")
}

// ─── gRPC generator (microservice + api-server) ──────────────────────────────

func TestGRPC_MicroserviceRequiredFiles(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)

	hasEntry(t, entries, "svc/README.md")
	hasEntry(t, entries, "svc/go.mod")
	hasEntry(t, entries, "svc/cmd/svc/main.go")
	hasEntry(t, entries, "svc/internal/server/server.go")
	hasEntry(t, entries, "svc/proto/svc.proto")
	hasEntry(t, entries, "svc/buf.yaml")
	hasEntry(t, entries, "svc/buf.gen.yaml")
	hasEntry(t, entries, "svc/Makefile")
	hasEntry(t, entries, "svc/.gitignore")
}

func TestGRPC_APIServerRequiredFiles(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "api-server",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/api",
		Name:        "api",
	}
	buf, err := GeneratorRegistry["api-server"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)

	hasEntry(t, entries, "api/README.md")
	hasEntry(t, entries, "api/go.mod")
	hasEntry(t, entries, "api/cmd/api/main.go")
	hasEntry(t, entries, "api/internal/server/server.go")
	hasEntry(t, entries, "api/proto/api.proto")
	hasEntry(t, entries, "api/buf.yaml")
	hasEntry(t, entries, "api/buf.gen.yaml")
	hasEntry(t, entries, "api/Makefile")
	hasEntry(t, entries, "api/.gitignore")
}

func TestGRPC_MainImportsGRPCDeps(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	main := hasEntry(t, entries, "svc/cmd/svc/main.go")
	assert.Contains(t, main, "google.golang.org/grpc")
	assert.Contains(t, main, "grpc_health_v1")
	assert.Contains(t, main, "reflection")
	assert.Contains(t, main, "package main")
}

func TestGRPC_ServerImplementsHealthCheck(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	srv := hasEntry(t, entries, "svc/internal/server/server.go")
	assert.Contains(t, srv, "HealthServer")
	assert.Contains(t, srv, "grpc_health_v1.UnimplementedHealthServer")
	assert.Contains(t, srv, "Check")
	assert.Contains(t, srv, "Watch")
}

func TestGRPC_ProtoFileContainsGreeterService(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	proto := hasEntry(t, entries, "svc/proto/svc.proto")
	assert.Contains(t, proto, `syntax = "proto3"`)
	assert.Contains(t, proto, "GreeterService")
	assert.Contains(t, proto, "SayHello")
	assert.Contains(t, proto, "go_package")
}

func TestGRPC_GoModContainsGRPCDeps(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	gomod := hasEntry(t, entries, "svc/go.mod")
	assert.Contains(t, gomod, "google.golang.org/grpc")
	assert.Contains(t, gomod, "google.golang.org/protobuf")
}

func TestGRPC_MakefileContainsProtoTarget(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	makefile := hasEntry(t, entries, "svc/Makefile")
	assert.Contains(t, makefile, "proto:")
	assert.Contains(t, makefile, "buf generate")
}

func TestGRPC_BufYamlVersion(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	bufYaml := hasEntry(t, entries, "svc/buf.yaml")
	assert.Contains(t, bufYaml, "version: v2")
	assert.Contains(t, bufYaml, "path: proto")
}

func TestGRPC_BufGenYamlPlugins(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	bufGenYaml := hasEntry(t, entries, "svc/buf.gen.yaml")
	assert.Contains(t, bufGenYaml, "buf.build/protocolbuffers/go")
	assert.Contains(t, bufGenYaml, "buf.build/grpc/go")
	assert.Contains(t, bufGenYaml, "out: gen")
}

func TestGRPC_GitignoreExcludesGenDir(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	gitignore := hasEntry(t, entries, "svc/.gitignore")
	assert.Contains(t, gitignore, "gen/")
}

func TestGRPC_ReadmeExplainsBufGenerate(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
		Description: "my grpc service",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	readme := hasEntry(t, entries, "svc/README.md")
	assert.Contains(t, readme, "buf generate")
	assert.Contains(t, readme, "svc")
	assert.Contains(t, readme, "my grpc service")
}

func TestGRPC_WithDockerSupport(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType:   "microservice",
		GoVersion:     "1.24.6",
		Framework:     "grpc",
		ModuleName:    "github.com/acme/svc",
		Name:          "svc",
		DockerSupport: true,
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "svc/Dockerfile")
}

func TestGRPC_WithoutDockerSupport(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType:   "microservice",
		GoVersion:     "1.24.6",
		Framework:     "grpc",
		ModuleName:    "github.com/acme/svc",
		Name:          "svc",
		DockerSupport: false,
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	_, ok := entries["svc/Dockerfile"]
	assert.False(t, ok, "Dockerfile should not be included when DockerSupport=false")
}

func TestGRPC_WithCacheAddon(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/svc",
		Name:        "svc",
		Addons:      map[string][]string{"cache": {"redis"}},
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "svc/internal/cache/cache.go")
}

func TestGRPC_HyphenatedNameProtoPackage(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "grpc",
		ModuleName:  "github.com/acme/my-service",
		Name:        "my-service",
	}
	buf, err := GeneratorRegistry["microservice"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	proto := hasEntry(t, entries, "my-service/proto/my-service.proto")
	// Proto package name must use underscores, not hyphens
	assert.Contains(t, proto, "package my_service.v1")
}

// ─── AIAgentGenerator ────────────────────────────────────────────────────────

var aiAgentFrameworks = []string{"langchaingo", "openai", "gemini", "ollama"}

func TestAIAgentGenerator_RequiredFiles(t *testing.T) {
	for _, fw := range aiAgentFrameworks {
		t.Run(fw, func(t *testing.T) {
			req := CreateProjectRequest{
				ProjectType: "ai-agent",
				GoVersion:   "1.24.6",
				Framework:   fw,
				ModuleName:  "github.com/acme/agent",
				Name:        "agent",
			}
			buf, err := GeneratorRegistry["ai-agent"].Generate(context.Background(), req)
			require.NoError(t, err)
			entries := zipEntries(t, buf)

			hasEntry(t, entries, "agent/README.md")
			hasEntry(t, entries, "agent/go.mod")
			hasEntry(t, entries, "agent/main.go")
			hasEntry(t, entries, "agent/agent/agent.go")
			hasEntry(t, entries, "agent/tools/tools.go")
			hasEntry(t, entries, "agent/llm/client.go")
			hasEntry(t, entries, "agent/Makefile")
			hasEntry(t, entries, "agent/.gitignore")
			hasEntry(t, entries, "agent/agent/agent_test.go")
		})
	}
}

func TestAIAgentGenerator_ReadmeContainsFrameworkAndName(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "ai-agent",
		GoVersion:   "1.24.6",
		Framework:   "openai",
		ModuleName:  "github.com/acme/myagent",
		Name:        "myagent",
		Description: "an intelligent assistant",
	}
	buf, err := GeneratorRegistry["ai-agent"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	readme := hasEntry(t, entries, "myagent/README.md")
	assert.Contains(t, readme, "myagent")
	assert.Contains(t, readme, "an intelligent assistant")
	assert.Contains(t, readme, "openai")
}

func TestAIAgentGenerator_WithDocker(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType:   "ai-agent",
		GoVersion:     "1.24.6",
		Framework:     "openai",
		ModuleName:    "github.com/acme/agent",
		Name:          "agent",
		DockerSupport: true,
	}
	buf, err := GeneratorRegistry["ai-agent"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "agent/Dockerfile")
}

func TestAIAgentGenerator_WithoutDocker(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType:   "ai-agent",
		GoVersion:     "1.24.6",
		Framework:     "openai",
		ModuleName:    "github.com/acme/agent",
		Name:          "agent",
		DockerSupport: false,
	}
	buf, err := GeneratorRegistry["ai-agent"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	_, hasDocker := entries["agent/Dockerfile"]
	assert.False(t, hasDocker, "Dockerfile should not be present when DockerSupport=false")
}

func TestAIAgentGenerator_OpenAIEnvExample(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "ai-agent",
		GoVersion:   "1.24.6",
		Framework:   "openai",
		ModuleName:  "github.com/acme/agent",
		Name:        "agent",
	}
	buf, err := GeneratorRegistry["ai-agent"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	env := hasEntry(t, entries, "agent/.env.example")
	assert.Contains(t, env, "OPENAI_API_KEY")
}

func TestAIAgentGenerator_GeminiEnvExample(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "ai-agent",
		GoVersion:   "1.24.6",
		Framework:   "gemini",
		ModuleName:  "github.com/acme/agent",
		Name:        "agent",
	}
	buf, err := GeneratorRegistry["ai-agent"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	env := hasEntry(t, entries, "agent/.env.example")
	assert.Contains(t, env, "GEMINI_API_KEY")
}

func TestAIAgentGenerator_OllamaNoEnvExample(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "ai-agent",
		GoVersion:   "1.24.6",
		Framework:   "ollama",
		ModuleName:  "github.com/acme/agent",
		Name:        "agent",
	}
	buf, err := GeneratorRegistry["ai-agent"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	_, hasEnv := entries["agent/.env.example"]
	assert.False(t, hasEnv, ".env.example should not be present for ollama (no API key needed)")
}

func TestAIAgentGenerator_GoModContainsFrameworkDep(t *testing.T) {
	cases := []struct {
		framework string
		dep       string
	}{
		{"langchaingo", "github.com/tmc/langchaingo"},
		{"openai", "github.com/openai/openai-go"},
		{"gemini", "github.com/google/generative-ai-go"},
		{"ollama", "github.com/ollama/ollama"},
	}
	for _, tc := range cases {
		t.Run(tc.framework, func(t *testing.T) {
			req := CreateProjectRequest{
				ProjectType: "ai-agent",
				GoVersion:   "1.24.6",
				Framework:   tc.framework,
				ModuleName:  "github.com/acme/agent",
				Name:        "agent",
			}
			buf, err := GeneratorRegistry["ai-agent"].Generate(context.Background(), req)
			require.NoError(t, err)
			entries := zipEntries(t, buf)
			gomod := hasEntry(t, entries, "agent/go.mod")
			assert.Contains(t, gomod, tc.dep)
		})
	}
}

func TestAIAgentGenerator_DockerfileUsesBuildRoot(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType:   "ai-agent",
		GoVersion:     "1.24.6",
		Framework:     "openai",
		ModuleName:    "github.com/acme/myagent",
		Name:          "myagent",
		DockerSupport: true,
	}
	buf, err := GeneratorRegistry["ai-agent"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	dockerfile := hasEntry(t, entries, "myagent/Dockerfile")
	// ai-agent uses root "." as build target
	assert.Contains(t, dockerfile, "go build -o myagent .")
}

func TestAIAgentGenerator_WithVectorStoreAddon(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "ai-agent",
		GoVersion:   "1.24.6",
		Framework:   "openai",
		ModuleName:  "github.com/acme/agent",
		Name:        "agent",
		Addons:      map[string][]string{"vectorstore": {"pgvector"}},
	}
	buf, err := GeneratorRegistry["ai-agent"].Generate(context.Background(), req)
	require.NoError(t, err)
	entries := zipEntries(t, buf)
	hasEntry(t, entries, "agent/internal/vectorstore/store.go")
}
