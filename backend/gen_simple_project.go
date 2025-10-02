package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"log"
)

// myproject/
// ├── cmd/
// │   └── myproject/
// │       └── main.go
// ├── internal/           # Business logic
// │   └── service.go
// ├── go.mod
// └── README.md
func GenerateSimpleProjecet(request CreateProjectRequest) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	folderName := request.Name
	if folderName == "" {
		folderName = "myproject"
	}

	// 1. README.md
	readmeContent := fmt.Sprintf("# %s\n\n%s", folderName, request.Description)
	readmeFile, err := zipWriter.Create(fmt.Sprintf("%s/README.md", folderName))
	if err != nil {
		log.Printf("[ERROR] Failed to create README.md in zip: %v", err)
		return nil, errors.New("failed to create README.md in zip")
	}
	_, err = readmeFile.Write([]byte(readmeContent))
	if err != nil {
		log.Printf("[ERROR] Failed to write README.md: %v", err)
		return nil, errors.New("failed to write README.md")
	}

	// 2. go.mod (minimal content)
	gomodContent, err := GenerateGoModV2(request)
	if err != nil {
		log.Printf("[ERROR] Failed to generate go.mod content: %v", err)
		return nil, errors.New("failed to generate go.mod content")
	}
	gomodFile, err := zipWriter.Create(fmt.Sprintf("%s/go.mod", folderName))
	if err != nil {
		log.Printf("[ERROR] Failed to create go.mod in zip: %v", err)
		return nil, errors.New("failed to create go.mod in zip")
	}
	_, err = gomodFile.Write([]byte(gomodContent))
	if err != nil {
		log.Printf("[ERROR] Failed to write go.mod: %v", err)
		return nil, errors.New("failed to write go.mod")
	}

	// 3. cmd/myproject/main.go using GenerateMain
	mainGoPath := fmt.Sprintf("%s/cmd/%s/main.go", folderName, folderName)
	mainGoFile, err := zipWriter.Create(mainGoPath)
	if err != nil {
		log.Printf("[ERROR] Failed to create main.go in zip: %v", err)
		return nil, errors.New("failed to create main.go in zip")
	}
	mainGoContent, err := GenerateMainContent()
	if err != nil {
		log.Printf("[ERROR] Failed to generate main.go content: %v", err)
		return nil, errors.New("failed to generate main.go content")
	}
	_, err = mainGoFile.Write(mainGoContent)
	if err != nil {
		log.Printf("[ERROR] Failed to write main.go: %v", err)
		return nil, errors.New("failed to write main.go")
	}

	// check addons and include if any
	if len(request.Addons) > 0 {
		for addonType, addons := range request.Addons {
			if addonType == "logging" {
			}
			if addonType == "database" {
				dbPath := fmt.Sprintf("%s/internal/database/database.go", folderName)
				dbFile, err := zipWriter.Create(dbPath)
				if err != nil {
					log.Printf("[ERROR] Failed to create database.go in zip: %v", err)
					return nil, errors.New("failed to create database.go in zip")
				}
				dbContent, err := GenerateDatabaseAddon(addons)
				if err != nil {
					log.Printf("[ERROR] Failed to generate database.go content: %v", err)
					return nil, errors.New("failed to generate database.go content")
				}
				_, err = dbFile.Write(dbContent)
				if err != nil {
					log.Printf("[ERROR] Failed to write database.go: %v", err)
					return nil, errors.New("failed to write database.go")
				}
			}
			if addonType == "cache" {
				cachePath := fmt.Sprintf("%s/internal/cache/cache.go", folderName)
				cacheFile, err := zipWriter.Create(cachePath)
				if err != nil {
					log.Printf("[ERROR] Failed to create cache.go in zip: %v", err)
					return nil, errors.New("failed to create cache.go in zip")
				}
				cacheContent, err := GenerateCacheAddon(addons)
				if err != nil {
					log.Printf("[ERROR] Failed to generate cache.go content: %v", err)
					return nil, errors.New("failed to generate cache.go content")
				}
				_, err = cacheFile.Write(cacheContent)
				if err != nil {
					log.Printf("[ERROR] Failed to write cache.go: %v", err)
					return nil, errors.New("failed to write cache.go")
				}
			}
		}
	}

	// 4. internal/service.go
	serviceGoContent := "package internal\n\n// Business logic goes here\nfunc Service() string {\n    return \"Service logic\"\n}\n"
	serviceGoPath := fmt.Sprintf("%s/internal/service.go", folderName)
	serviceGoFile, err := zipWriter.Create(serviceGoPath)
	if err != nil {
		log.Printf("[ERROR] Failed to create service.go in zip: %v", err)
		return nil, errors.New("failed to create service.go in zip")
	}
	_, err = serviceGoFile.Write([]byte(serviceGoContent))
	if err != nil {
		log.Printf("[ERROR] Failed to write service.go: %v", err)
		return nil, errors.New("failed to write service.go")
	}

	err = zipWriter.Close()
	if err != nil {
		log.Printf("[ERROR] Failed to close zip writer: %v", err)
		return nil, errors.New("failed to finalize zip file")
	}

	return buf, nil
}
