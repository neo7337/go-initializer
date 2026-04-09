package generator

import (
	"archive/zip"
	"bytes"
	"context"
)

// Generator is implemented by each project-type generator.
type Generator interface {
	Generate(ctx context.Context, request CreateProjectRequest) (*bytes.Buffer, error)
}

// AddonGenerator writes addon-specific files directly into the project zip.
type AddonGenerator interface {
	Generate(folderName string, addons []string, zw *zip.Writer) error
}

// GeneratorRegistry maps project types to their Generator implementations.
// To add a new project type, register it here and add a Generator implementation.
var GeneratorRegistry = map[string]Generator{
	"simple-project": &SimpleProjectGenerator{},
	"microservice":   &MicroserviceGenerator{},
	"cli-app":        &CLIAppGenerator{},
	"api-server":     &APIServerGenerator{},
	"ai-agent":       &AIAgentGenerator{},
}

// addonRegistry maps addon category names to their AddonGenerator implementations.
// To add a new addon category, register it here and add an AddonGenerator implementation.
var addonRegistry = map[string]AddonGenerator{
	"cache":       &CacheAddonGen{},
	"database":    &DatabaseAddonGen{},
	"ai":          &AIAddonGen{},
	"vectorstore": &VectorStoreAddonGen{},
}

var SupportedProjectTypesMap = map[string]bool{
	"microservice":   true,
	"simple-project": true,
	"cli-app":        true,
	"api-server":     true,
	"ai-agent":       true,
}

var SupportedProjectTypesLabelsMap = map[string]string{
	"microservice":   "Microservice",
	"simple-project": "Simple Project",
	"cli-app":        "CLI Application",
	"api-server":     "API Server",
	"ai-agent":       "AI Agent",
}

var SupportedGoVersionsMap = map[string]bool{
	"1.25.0":  true,
	"1.24.6":  true,
	"1.23.12": true,
}

var SupportedFrameworksMap = map[string]map[string]bool{
	"microservice": {
		"golly": true,
		"gin":   true,
		"echo":  true,
		"fiber": true,
		"gokit": true,
	},
	"cli-app": {
		"golly":   true,
		"cobra":   true,
		"urfave":  true,
		"kingpin": true,
	},
	"api-server": {
		"golly": true,
		"gin":   true,
		"echo":  true,
		"fiber": true,
		"chi":   true,
	},
	"simple-project": {
		"golly": true,
	},
	"ai-agent": {
		"langchaingo": true,
		"openai":      true,
		"gemini":      true,
		"ollama":      true,
	},
}

var SupportedAddonsMap = map[string]map[string]bool{
	"cache": {
		"redis":     true,
		"memcached": true,
	},
	"database": {
		"gorm": true,
		"ent":  true,
	},
	"other": {
		"zap":    true,
		"logrus": true,
		"cobra":  true,
	},
	"ai": {
		"openai":      true,
		"langchaingo": true,
		"gemini":      true,
		"ollama":      true,
	},
	"vectorstore": {
		"pgvector": true,
		"chromem":  true,
		"qdrant":   true,
	},
}

var DependencyMap = map[string][]string{
	"golly":       {"oss.nandlabs.io/golly v1.2.9"},
	"gin":         {"github.com/gin-gonic/gin v1.11.0"},
	"echo":        {"github.com/labstack/echo/v4 v4.14.0"},
	"fiber":       {"github.com/gofiber/fiber/v2 v2.52.10"},
	"chi":         {"github.com/go-chi/chi/v5 v5.2.1"},
	"gorm":        {"gorm.io/gorm v1.31.1", "gorm.io/driver/postgres v1.6.0"},
	"ent":         {"entgo.io/ent/cmd/ent v0.14.5"},
	"sqlx":        {"github.com/jmoiron/sqlx v1.4.0", "github.com/lib/pq v1.10.9"},
	"zap":         {"go.uber.org/zap v1.27.1"},
	"logrus":      {"github.com/sirupsen/logrus v1.9.3"},
	"zerolog":     {"github.com/rs/zerolog v1.34.0"},
	"viper":       {"github.com/spf13/viper v1.21.0"},
	"cobra":       {"github.com/spf13/cobra v1.10.2", "github.com/spf13/pflag v1.0.10"},
	"urfave":      {"github.com/urfave/cli/v2 v2.27.6"},
	"kingpin":     {"gopkg.in/alecthomas/kingpin.v2 v2.2.6"},
	"testify":     {"github.com/stretchr/testify v1.11.1"},
	"httptest":    {"net/http/httptest"},
	"redis":       {"github.com/redis/go-redis/v9 v9.17.2"},
	"memcached":   {"github.com/bradfitz/gomemcache/memcache v0.0.0-20250403215159-8d39553ac7cf"},
	"langchaingo": {"github.com/tmc/langchaingo v0.1.13"},
	"openai":      {"github.com/openai/openai-go v0.1.0-beta.10"},
	"gemini":      {"github.com/google/generative-ai-go v0.19.0", "google.golang.org/api v0.228.0"},
	"ollama":      {"github.com/ollama/ollama v0.6.5"},
	"pgvector":    {"github.com/pgvector/pgvector-go v0.2.2"},
	"chromem":     {"github.com/philippgille/chromem-go v0.7.0"},
	"qdrant":      {"github.com/qdrant/go-client v1.13.0"},
}
