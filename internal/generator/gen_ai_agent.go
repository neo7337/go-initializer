package generator

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
)

// AIAgentGenerator implements Generator for the "ai-agent" project type.
//
// Generated layout:
//
//	<name>/
//	├── main.go                    # calls agent.Run()
//	├── agent/
//	│   └── agent.go               # ReAct-style LLM loop
//	├── tools/
//	│   └── tools.go               # sample tool stubs (function calling)
//	├── llm/
//	│   └── client.go              # LLM provider initialisation
//	├── internal/
//	│   ├── logger/logger.go       # if "other" logging addon selected
//	│   └── vectorstore/store.go   # if "vectorstore" addon selected
//	├── go.mod
//	├── .gitignore
//	├── Makefile
//	├── README.md
//	└── Dockerfile                 # if dockerSupport is true
type AIAgentGenerator struct{}

func (g *AIAgentGenerator) Generate(request CreateProjectRequest) (*bytes.Buffer, error) {
	if request.Name == "" {
		request.Name = "myagent"
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

	// main.go — thin entrypoint that calls agent.Run()
	if err := addToZip(zipWriter, fmt.Sprintf("%s/main.go", folderName), generateAIAgentMain(request.ModuleName)); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// agent/agent.go — ReAct-style LLM loop
	if err := addToZip(zipWriter, fmt.Sprintf("%s/agent/agent.go", folderName), GenerateAgentContent(request.Framework, folderName)); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// tools/tools.go — sample tool stubs
	if err := addToZip(zipWriter, fmt.Sprintf("%s/tools/tools.go", folderName), GenerateToolsContent(request.Framework)); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// llm/client.go — LLM provider initialisation
	if err := addToZip(zipWriter, fmt.Sprintf("%s/llm/client.go", folderName), GenerateLLMClient(request.Framework)); err != nil {
		log.Printf("[ERROR] %v", err)
		return nil, err
	}

	// Addons
	for addonType, addons := range request.Addons {
		if len(addons) == 0 {
			continue
		}
		switch addonType {
		case "other":
			// Logging addon — zap or logrus
			loggerContent, err := GenerateLoggingAddon(addons)
			if err != nil {
				log.Printf("[WARN] Skipping logging addon: %v", err)
				continue
			}
			if err := addToZip(zipWriter, fmt.Sprintf("%s/internal/logger/logger.go", folderName), loggerContent); err != nil {
				log.Printf("[ERROR] %v", err)
				return nil, err
			}
		case "vectorstore":
			// Vectorstore addon — write the first selected store to internal/vectorstore/store.go
			for _, store := range addons {
				if err := addToZip(zipWriter, fmt.Sprintf("%s/internal/vectorstore/store.go", folderName), GenerateVectorStoreContent(store)); err != nil {
					log.Printf("[ERROR] %v", err)
					return nil, err
				}
				// Only one store file is generated; subsequent selections in the same
				// category would overwrite — stop after the first.
				break
			}
		default:
			gen, ok := addonRegistry[addonType]
			if !ok {
				continue
			}
			if err := gen.Generate(folderName, addons, zipWriter); err != nil {
				log.Printf("[ERROR] Failed to generate %s addon: %v", addonType, err)
				return nil, fmt.Errorf("failed to generate %s addon: %w", addonType, err)
			}
		}
	}

	// Makefile — main package is at the project root "."
	if err := addToZip(zipWriter, fmt.Sprintf("%s/Makefile", folderName), GenerateMakefile(folderName, ".")); err != nil {
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

// generateAIAgentMain returns a thin main.go that delegates to agent.Run().
func generateAIAgentMain(moduleName string) []byte {
	return []byte(fmt.Sprintf(`// Code generated by go-initializer. DO NOT EDIT.
package main

import "%s/agent"

func main() {
	agent.Run()
}
`, moduleName))
}
