package main

type CreateProjectRequest struct {
	ProjectType   string              `json:"projectType"`
	GoVersion     string              `json:"goVersion"`
	Framework     string              `json:"framework"`
	ModuleName    string              `json:"moduleName"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Addons        map[string][]string `json:"selectedAddons,omitempty"` // New field for addons
	DockerSupport bool                `json:"dockerSupport"`            // New field for Docker support
}

var SupportedProjectTypesMap = map[string]bool{
	"microservice":   true,
	"simple-project": true,
	"cli-app":        true,
	"api-server":     true,
}

var SupportedProjectTypesLabelsMap = map[string]string{
	"microservice":   "Microservice",
	"simple-project": "Simple Project",
	"cli-app":        "CLI Application",
	"api-server":     "API Server",
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
}

var DependencyMap = map[string][]string{
	"golly":     {"oss.nandlabs.io/golly v1.2.9"},
	"gin":       {"github.com/gin-gonic/gin v1.11.0"},
	"echo":      {"github.com/labstack/echo/v4 v4.14.0"},
	"fiber":     {"github.com/gofiber/fiber/v2 v2.52.10"},
	"gorm":      {"gorm.io/gorm v1.31.1", "gorm.io/driver/postgres v1.6.0"},
	"ent":       {"entgo.io/ent/cmd/ent v0.14.5"},
	"sqlx":      {"github.com/jmoiron/sqlx v1.4.0", "github.com/lib/pq v1.10.9"},
	"zap":       {"go.uber.org/zap v1.27.1"},
	"logrus":    {"github.com/sirupsen/logrus v1.9.3"},
	"zerolog":   {"github.com/rs/zerolog v1.34.0"},
	"viper":     {"github.com/spf13/viper v1.21.0"},
	"cobra":     {"github.com/spf13/cobra v1.10.2", "github.com/spf13/pflag v1.0.10"},
	"testify":   {"github.com/stretchr/testify v1.11.1"},
	"httptest":  {"net/http/httptest"},
	"redis":     {"github.com/redis/go-redis/v9 v9.17.2"},
	"memcached": {"github.com/bradfitz/gomemcache/memcache v0.0.0-20250403215159-8d39553ac7cf"},
}
