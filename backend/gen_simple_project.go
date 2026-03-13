package main

import (
	"archive/zip"
	"bytes"
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

func (g *SimpleProjectGenerator) Generate(request CreateProjectRequest) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	folderName := request.Name
	if folderName == "" {
		folderName = "myproject"
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

	// cmd/<name>/main.go
	mainGoContent, err := GenerateMainContent()
	if err != nil {
		log.Printf("[ERROR] Failed to generate main.go: %v", err)
		return nil, fmt.Errorf("failed to generate main.go: %w", err)
	}
	if err := addToZip(zipWriter, fmt.Sprintf("%s/cmd/%s/main.go", folderName, folderName), mainGoContent); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// Addons — delegated to the addon registry; new addon types need only a
	// registration entry in registry.go, no changes here.
	for addonType, addons := range request.Addons {
		gen, ok := addonRegistry[addonType]
		if !ok || len(addons) == 0 {
			continue
		}
		if err := gen.Generate(folderName, addons, zipWriter); err != nil {
			log.Printf("[ERROR] Failed to generate %s addon: %v", addonType, err)
			return nil, fmt.Errorf("failed to generate %s addon: %w", addonType, err)
		}
	}

	// internal/service.go
	serviceContent := "package internal\n\n// Business logic goes here\nfunc Service() string {\n\treturn \"Service logic\"\n}\n"
	if err := addToZip(zipWriter, fmt.Sprintf("%s/internal/service.go", folderName), []byte(serviceContent)); err != nil {
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
