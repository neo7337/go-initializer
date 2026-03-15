package generator

type CreateProjectRequest struct {
	ProjectType   string              `json:"projectType"   validate:"required,oneof=microservice simple-project cli-app api-server"`
	GoVersion     string              `json:"goVersion"     validate:"required,oneof=1.25.0 1.24.6 1.23.12"`
	Framework     string              `json:"framework"     validate:"required"`
	ModuleName    string              `json:"moduleName"    validate:"required"`
	Name          string              `json:"name"          validate:"required"`
	Description   string              `json:"description"`
	Addons        map[string][]string `json:"selectedAddons,omitempty"`
	DockerSupport bool                `json:"dockerSupport"`
}
