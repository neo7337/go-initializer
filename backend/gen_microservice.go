package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
)

// MicroserviceGenerator implements Generator for the "microservice" project type.
//
// Generated layout:
//
//	<name>/
//	├── cmd/<name>/main.go        # framework-specific entrypoint
//	├── internal/
//	│   ├── handler/handler.go    # framework-specific handler
//	│   ├── router/router.go      # framework-specific router
//	│   ├── service/service.go    # stub service interface + impl
//	│   ├── cache/cache.go        # if cache addon selected
//	│   ├── database/database.go  # if database addon selected
//	│   └── logger/logger.go      # if zap/logrus addon selected
//	├── go.mod
//	├── .gitignore
//	├── Makefile
//	├── README.md
//	└── Dockerfile                # if dockerSupport is true
type MicroserviceGenerator struct{}

func (g *MicroserviceGenerator) Generate(request CreateProjectRequest) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	folderName := request.Name
	if folderName == "" {
		folderName = "myservice"
	}

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

	// cmd/<name>/main.go — framework-aware
	mainContent, err := GenerateMainContent(request.Framework)
	if err != nil {
		log.Printf("[ERROR] Failed to generate main.go: %v", err)
		return nil, fmt.Errorf("failed to generate main.go: %w", err)
	}
	if err := addToZip(zipWriter, fmt.Sprintf("%s/cmd/%s/main.go", folderName, folderName), mainContent); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// internal/handler/handler.go — framework-specific handler
	if err := addToZip(zipWriter, fmt.Sprintf("%s/internal/handler/handler.go", folderName), GenerateHandlerContent(request.Framework)); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// internal/router/router.go — framework-specific router
	if err := addToZip(zipWriter, fmt.Sprintf("%s/internal/router/router.go", folderName), GenerateRouterContent(request.Framework)); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// internal/service/service.go — named service stub
	if err := addToZip(zipWriter, fmt.Sprintf("%s/internal/service/service.go", folderName), GenerateServiceContent(folderName)); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// Addons — cache/database via registry; "other" (logging) handled inline
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
