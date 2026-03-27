package generator

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"log"
)

// SimpleProjectGenerator implements Generator for the "simple-project" project type.
//
// Generated layout:
//
//	<name>/
//	├── cmd/<name>/main.go
//	├── internal/service.go
//	├── go.mod
//	├── README.md
//	└── Dockerfile          (optional)
type SimpleProjectGenerator struct{}

func (g *SimpleProjectGenerator) Generate(ctx context.Context, request CreateProjectRequest) (*bytes.Buffer, error) {
	if request.Name == "" {
		request.Name = "myproject"
	}
	if err := ValidateProjectName(request.Name); err != nil {
		return nil, err
	}
	if err := ValidateModuleName(request.ModuleName); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	folderName := request.Name

	// README.md
	readmeContent := fmt.Sprintf("# %s\n\n%s", folderName, request.Description)
	if err := addToZip(zipWriter, fmt.Sprintf("%s/README.md", folderName), []byte(readmeContent)); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// go.mod
	gomodContent, err := GenerateGoModV2(request)
	if err != nil {
		log.Printf("[ERROR] Failed to generate go.mod: %v", err)
		return nil, fmt.Errorf("failed to generate go.mod: %w", err)
	}
	if err := addToZip(zipWriter, fmt.Sprintf("%s/go.mod", folderName), gomodContent); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// cmd/<name>/main.go
	mainGoContent, err := GenerateMainContent(request.Framework)
	if err != nil {
		log.Printf("[ERROR] Failed to generate main.go: %v", err)
		return nil, fmt.Errorf("failed to generate main.go: %w", err)
	}
	if err := addToZip(zipWriter, fmt.Sprintf("%s/cmd/%s/main.go", folderName, folderName), mainGoContent); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// Addons — cache/database delegated to the addon registry; "other" (logging)
	// is handled inline via GenerateLoggingAddon.
	for addonType, addons := range request.Addons {
		if len(addons) == 0 {
			continue
		}
		if addonType == "other" {
			loggerContent, err := GenerateLoggingAddon(addons)
			if err != nil {
				log.Printf("[WARN] Skipping logging addon: %v", err)
				continue
			}
			if err := addToZip(zipWriter, fmt.Sprintf("%s/internal/logger/logger.go", folderName), loggerContent); err != nil {
				log.Printf("[ERROR] %v", err)
				return nil, err
			}
			continue
		}
		gen, ok := addonRegistry[addonType]
		if !ok {
			continue
		}
		if err := gen.Generate(folderName, addons, zipWriter); err != nil {
			log.Printf("[ERROR] Failed to generate %s addon: %v", addonType, err)
			return nil, fmt.Errorf("failed to generate %s addon: %w", addonType, err)
		}
	}

	// internal/service.go
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := addToZip(zipWriter, fmt.Sprintf("%s/internal/service.go", folderName), GenerateServiceContent(folderName)); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// Makefile
	if err := addToZip(zipWriter, fmt.Sprintf("%s/Makefile", folderName), GenerateMakefile(folderName, fmt.Sprintf("./cmd/%s", folderName))); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// .gitignore
	if err := addToZip(zipWriter, fmt.Sprintf("%s/.gitignore", folderName), GenerateGitignore()); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// Dockerfile (optional)
	if request.DockerSupport {
		if err := addToZip(zipWriter, fmt.Sprintf("%s/Dockerfile", folderName), GenerateDockerfile(request)); err != nil {
			log.Printf("[ERROR] %v", err)
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		log.Printf("[ERROR] Failed to close zip writer: %v", err)
		return nil, fmt.Errorf("failed to finalize zip: %w", err)
	}
	return buf, nil
}

// myproject/
// ├── cmd/
// │   └── myproject/
// │       └── main.go
