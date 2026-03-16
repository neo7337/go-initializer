package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── GenerateGoModV2 ─────────────────────────────────────────────────────────

func TestGenerateGoModV2_BasicFields(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/myapp",
		Name:        "myapp",
	}
	out, err := GenerateGoModV2(req)
	require.NoError(t, err)
	body := string(out)
	assert.Contains(t, body, "module github.com/acme/myapp")
	assert.Contains(t, body, "go 1.24.6")
	assert.Contains(t, body, "github.com/gin-gonic/gin")
}

func TestGenerateGoModV2_UnsupportedProjectType(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "unknown-type",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/myapp",
		Name:        "myapp",
	}
	_, err := GenerateGoModV2(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported project type")
}

func TestGenerateGoModV2_UnsupportedFramework(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "rails",
		ModuleName:  "github.com/acme/myapp",
		Name:        "myapp",
	}
	_, err := GenerateGoModV2(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported framework")
}

func TestGenerateGoModV2_WithCacheAddon(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/myapp",
		Name:        "myapp",
		Addons:      map[string][]string{"cache": {"redis"}},
	}
	out, err := GenerateGoModV2(req)
	require.NoError(t, err)
	assert.Contains(t, string(out), "github.com/redis/go-redis")
}

func TestGenerateGoModV2_WithDatabaseAddon(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "microservice",
		GoVersion:   "1.24.6",
		Framework:   "gin",
		ModuleName:  "github.com/acme/myapp",
		Name:        "myapp",
		Addons:      map[string][]string{"database": {"gorm"}},
	}
	out, err := GenerateGoModV2(req)
	require.NoError(t, err)
	assert.Contains(t, string(out), "gorm.io/gorm")
}

func TestGenerateGoModV2_NoAddons(t *testing.T) {
	req := CreateProjectRequest{
		ProjectType: "simple-project",
		GoVersion:   "1.25.0",
		Framework:   "golly",
		ModuleName:  "github.com/acme/proj",
		Name:        "proj",
	}
	out, err := GenerateGoModV2(req)
	require.NoError(t, err)
	// golly has a dependency so there will be a require block, but no addon deps
	assert.NotContains(t, string(out), "gorm.io")
	assert.NotContains(t, string(out), "redis")
}

func TestGenerateGoModV2_AllSupportedFrameworksPerType(t *testing.T) {
	for projectType, frameworks := range SupportedFrameworksMap {
		for fw := range frameworks {
			t.Run(projectType+"/"+fw, func(t *testing.T) {
				req := CreateProjectRequest{
					ProjectType: projectType,
					GoVersion:   "1.24.6",
					Framework:   fw,
					ModuleName:  "github.com/acme/x",
					Name:        "x",
				}
				out, err := GenerateGoModV2(req)
				require.NoError(t, err)
				assert.Contains(t, string(out), "module github.com/acme/x")
			})
		}
	}
}

// ─── GenerateMainContent ─────────────────────────────────────────────────────

var mainContentCases = []struct {
	framework string
	contains  []string
}{
	{"gin", []string{"package main", "gin.Default()", `r.Run(":8080")`}},
	{"echo", []string{"package main", "echo.New()", `e.Start(":8080")`}},
	{"fiber", []string{"package main", "fiber.New()", `app.Listen(":8080")`}},
	{"chi", []string{"package main", "chi.NewRouter()", `http.ListenAndServe(":8080"`}},
	{"cobra", []string{"package main", "cobra.Command", "rootCmd.Execute()"}},
	{"gokit", []string{"package main", `http.ListenAndServe(":8080"`, `signal.Notify`}},
	{"golly", []string{"package main", "l3.Get()"}},
}

func TestGenerateMainContent_AllFrameworks(t *testing.T) {
	for _, tc := range mainContentCases {
		t.Run(tc.framework, func(t *testing.T) {
			out, err := GenerateMainContent(tc.framework)
			require.NoError(t, err)
			body := string(out)
			for _, want := range tc.contains {
				assert.Contains(t, body, want, "framework %s should contain %q", tc.framework, want)
			}
		})
	}
}

func TestGenerateMainContent_UnknownFrameworkDefaultsToGolly(t *testing.T) {
	out, err := GenerateMainContent("unknown")
	require.NoError(t, err)
	body := string(out)
	assert.Contains(t, body, "package main")
	assert.Contains(t, body, "l3.Get()")
}

// ─── GenerateHandlerContent ──────────────────────────────────────────────────

func TestGenerateHandlerContent_AllFrameworks(t *testing.T) {
	cases := []struct {
		framework string
		contains  string
	}{
		{"gin", "gin.Context"},
		{"echo", "echo.Context"},
		{"fiber", "fiber.Ctx"},
		{"chi", "http.ResponseWriter"},
		{"gokit", "MakeHealthEndpoint"},
		{"golly", "l3.Get()"},
	}
	for _, tc := range cases {
		t.Run(tc.framework, func(t *testing.T) {
			out := GenerateHandlerContent(tc.framework)
			assert.NotEmpty(t, out)
			body := string(out)
			assert.Contains(t, body, "package handler")
			assert.Contains(t, body, tc.contains)
		})
	}
}

func TestGenerateHandlerContent_UnknownDefaultsToGolly(t *testing.T) {
	out := GenerateHandlerContent("unknown")
	body := string(out)
	assert.Contains(t, body, "package handler")
	assert.Contains(t, body, "l3.Get()")
}

// ─── GenerateRouterContent ───────────────────────────────────────────────────

func TestGenerateRouterContent_AllFrameworks(t *testing.T) {
	cases := []struct {
		framework string
		contains  string
	}{
		{"gin", "*gin.Engine"},
		{"echo", "*echo.Echo"},
		{"fiber", "*fiber.App"},
		{"chi", "chi.NewRouter()"},
		{"gokit", "http.Handler"},
		{"golly", "l3.Get()"},
	}
	for _, tc := range cases {
		t.Run(tc.framework, func(t *testing.T) {
			out := GenerateRouterContent(tc.framework)
			assert.NotEmpty(t, out)
			body := string(out)
			assert.Contains(t, body, "package router")
			assert.Contains(t, body, tc.contains)
		})
	}
}

func TestGenerateRouterContent_UnknownDefaultsToGolly(t *testing.T) {
	out := GenerateRouterContent("unknown")
	body := string(out)
	assert.Contains(t, body, "package router")
	assert.Contains(t, body, "l3.Get()")
}

// ─── GenerateServiceContent ──────────────────────────────────────────────────

func TestGenerateServiceContent_InterfaceAndImpl(t *testing.T) {
	out := GenerateServiceContent("order")
	body := string(out)
	assert.Contains(t, body, "OrderService")
	assert.Contains(t, body, "orderService")
	assert.Contains(t, body, "NewOrderService")
	assert.Contains(t, body, "Hello()")
}

func TestGenerateServiceContent_UsesProvidedName(t *testing.T) {
	out := GenerateServiceContent("payment")
	body := string(out)
	assert.Contains(t, body, "PaymentService")
	assert.Contains(t, body, "NewPaymentService")
	assert.Contains(t, body, "Hello from payment")
}

// ─── GenerateGitignore ───────────────────────────────────────────────────────

func TestGenerateGitignore_ContainsExpectedPatterns(t *testing.T) {
	out := GenerateGitignore()
	body := string(out)
	assert.Contains(t, body, "vendor/")
	assert.Contains(t, body, ".env")
	assert.Contains(t, body, "*.exe")
	assert.Contains(t, body, ".DS_Store")
}

// ─── GenerateMakefile ────────────────────────────────────────────────────────

func TestGenerateMakefile_ContainsTargets(t *testing.T) {
	out := GenerateMakefile("myapp", "./cmd/myapp")
	body := string(out)
	assert.Contains(t, body, "BINARY_NAME=myapp")
	assert.Contains(t, body, "MAIN_PKG=./cmd/myapp")
	assert.Contains(t, body, "build:")
	assert.Contains(t, body, "run:")
	assert.Contains(t, body, "test:")
	assert.Contains(t, body, "tidy:")
}

// ─── GenerateDockerfile ──────────────────────────────────────────────────────

func TestGenerateDockerfile_MultistageContent(t *testing.T) {
	req := CreateProjectRequest{Name: "myapp"}
	out := GenerateDockerfile(req)
	body := string(out)
	assert.Contains(t, body, "FROM golang:alpine AS builder")
	assert.Contains(t, body, "FROM alpine:latest")
	assert.Contains(t, body, "EXPOSE 8080")
}

// ─── GenerateLoggingAddon ────────────────────────────────────────────────────

func TestGenerateLoggingAddon_Zap(t *testing.T) {
	out, err := GenerateLoggingAddon([]string{"zap"})
	require.NoError(t, err)
	body := string(out)
	assert.Contains(t, body, "go.uber.org/zap")
	assert.Contains(t, body, "zap.NewProduction()")
}

func TestGenerateLoggingAddon_Logrus(t *testing.T) {
	out, err := GenerateLoggingAddon([]string{"logrus"})
	require.NoError(t, err)
	body := string(out)
	assert.Contains(t, body, "github.com/sirupsen/logrus")
	assert.Contains(t, body, "logrus.New()")
}

func TestGenerateLoggingAddon_UnknownAddon(t *testing.T) {
	_, err := GenerateLoggingAddon([]string{"unknown-logger"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no supported logging addon")
}

func TestGenerateLoggingAddon_ZapTakesPriority(t *testing.T) {
	out, err := GenerateLoggingAddon([]string{"zap", "unknown"})
	require.NoError(t, err)
	assert.Contains(t, string(out), "zap")
}

// ─── GenerateCacheAddon ──────────────────────────────────────────────────────

func TestGenerateCacheAddon_Redis(t *testing.T) {
	out, err := GenerateCacheAddon([]string{"redis"})
	require.NoError(t, err)
	require.NotNil(t, out)
	body := string(out)
	assert.Contains(t, body, "NewRedisClient")
	assert.Contains(t, body, "github.com/redis/go-redis")
}

func TestGenerateCacheAddon_Memcached(t *testing.T) {
	out, err := GenerateCacheAddon([]string{"memcached"})
	require.NoError(t, err)
	require.NotNil(t, out)
	body := string(out)
	assert.Contains(t, body, "NewMemcacheClient")
	assert.Contains(t, body, "gomemcache")
}

func TestGenerateCacheAddon_Unknown(t *testing.T) {
	_, err := GenerateCacheAddon([]string{"unknown-cache"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported cache addon")
}

func TestGenerateCacheAddon_Empty(t *testing.T) {
	out, err := GenerateCacheAddon([]string{})
	require.NoError(t, err)
	assert.Nil(t, out)
}

// ─── GenerateDatabaseAddon ───────────────────────────────────────────────────

func TestGenerateDatabaseAddon_GORM(t *testing.T) {
	out, err := GenerateDatabaseAddon([]string{"gorm"})
	require.NoError(t, err)
	require.NotNil(t, out)
	body := string(out)
	assert.Contains(t, body, "NewGormDB")
	assert.Contains(t, body, "gorm.io/driver/postgres")
}

func TestGenerateDatabaseAddon_Ent(t *testing.T) {
	out, err := GenerateDatabaseAddon([]string{"ent"})
	require.NoError(t, err)
	require.NotNil(t, out)
	body := string(out)
	assert.Contains(t, body, "NewEntClient")
	assert.Contains(t, body, "DATABASE_URL")
}

func TestGenerateDatabaseAddon_Unknown(t *testing.T) {
	_, err := GenerateDatabaseAddon([]string{"mysql-unsupported"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database addon")
}

func TestGenerateDatabaseAddon_Empty(t *testing.T) {
	out, err := GenerateDatabaseAddon([]string{})
	require.NoError(t, err)
	assert.Nil(t, out)
}

// ─── GenerateRootCmd ─────────────────────────────────────────────────────────

func TestGenerateRootCmd_AllFrameworks(t *testing.T) {
	cases := []struct {
		framework string
		contains  string
	}{
		{"cobra", "cobra.Command"},
		{"urfave", "cli.App"},
		{"kingpin", "kingpin.New"},
		{"golly", "l3.Get()"},
	}
	for _, tc := range cases {
		t.Run(tc.framework, func(t *testing.T) {
			out := GenerateRootCmd(tc.framework, "myapp")
			body := string(out)
			assert.NotEmpty(t, body)
			assert.Contains(t, body, "package cmd")
			assert.Contains(t, body, tc.contains)
			assert.Contains(t, body, "Execute")
		})
	}
}

func TestGenerateRootCmd_UsesBinaryName(t *testing.T) {
	out := GenerateRootCmd("cobra", "myservice")
	assert.Contains(t, string(out), `"myservice"`)
}

// ─── GenerateSubCmd ──────────────────────────────────────────────────────────

func TestGenerateSubCmd_AllFrameworks(t *testing.T) {
	cases := []struct {
		framework string
		contains  string
	}{
		{"cobra", "cobra.Command"},
		{"urfave", "cli.Command"},
		{"kingpin", `App.Command(`},
		{"golly", "l3.Get()"},
	}
	for _, tc := range cases {
		t.Run(tc.framework, func(t *testing.T) {
			out := GenerateSubCmd(tc.framework, "myapp")
			body := string(out)
			assert.NotEmpty(t, body)
			assert.Contains(t, body, "package cmd")
			assert.Contains(t, body, tc.contains)
		})
	}
}

func TestGenerateSubCmd_UsesName(t *testing.T) {
	out := GenerateSubCmd("cobra", "deployer")
	assert.Contains(t, string(out), "deployer")
}

// ─── helpers ─────────────────────────────────────────────────────────────────

// containsAllOf asserts that str contains every item in wants.
func containsAllOf(t *testing.T, str string, wants ...string) {
	t.Helper()
	for _, w := range wants {
		assert.True(t, strings.Contains(str, w), "expected %q to contain %q", str, w)
	}
}
