package main

type CreateProjectRequest struct {
	ProjectType string `json:"projectType"`
	GoVersion   string `json:"goVersion"`
	Framework   string `json:"framework"`
	ModuleName  string `json:"moduleName"`
	Name        string `json:"name"`
	Description string `json:"description"`
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
	"golly":    {"github.com/nandlabs/golly"},
	"gin":      {"github.com/gin-gonic/gin"},
	"echo":     {"github.com/labstack/echo/v4"},
	"fiber":    {"github.com/gofiber/fiber/v2"},
	"gorm":     {"gorm.io/gorm", "gorm.io/driver/postgres"},
	"ent":      {"entgo.io/ent/cmd/ent"},
	"sqlx":     {"github.com/jmoiron/sqlx", "github.com/lib/pq"},
	"zap":      {"go.uber.org/zap"},
	"logrus":   {"github.com/sirupsen/logrus"},
	"zerolog":  {"github.com/rs/zerolog"},
	"viper":    {"github.com/spf13/viper"},
	"cobra":    {"github.com/spf13/cobra", "github.com/spf13/pflag"},
	"testify":  {"github.com/stretchr/testify"},
	"httptest": {"net/http/httptest"},
}
