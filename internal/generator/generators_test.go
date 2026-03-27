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
	for _, typ := range []string{"simple-project", "microservice", "cli-app"} {
		_, ok := GeneratorRegistry[typ]
		assert.True(t, ok, "expected %q to be in GeneratorRegistry", typ)
	}
}

func TestGeneratorRegistry_APIServerNotRegistered(t *testing.T) {
	// api-server is in the validation allowlist but not yet implemented
	_, ok := GeneratorRegistry["api-server"]
	assert.False(t, ok, "api-server should not be in GeneratorRegistry yet")
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
