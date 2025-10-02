package main

import (
	"fmt"
)

func GenerateGoMod(request CreateProjectRequest) string {
	requiredImports := make([]string, 0)

	// Validate project type
	if !SupportedProjectTypesMap[request.ProjectType] {
		return fmt.Sprintf("// Unsupported project type: %s\n", request.ProjectType)
	}

	// Validate framework for the project type
	frameworks, ok := SupportedFrameworksMap[request.ProjectType]
	if !ok || !frameworks[request.Framework] {
		return fmt.Sprintf("// Unsupported framework '%s' for project type '%s'\n", request.Framework, request.ProjectType)
	}

	// Addons handling
	for addonType, addons := range request.Addons {
		if supportedAddons, exists := SupportedAddonsMap[addonType]; exists {
			for _, addon := range addons {
				if supportedAddons[addon] {
					if deps, ok := DependencyMap[addon]; ok {
						requiredImports = append(requiredImports, deps...)
					}
				}
			}
		}
	}

	// Get dependencies for the framework
	if deps, ok := DependencyMap[request.Framework]; ok {
		requiredImports = append(requiredImports, deps...)
	} else {
		requiredImports = []string{}
	}

	goMod := fmt.Sprintf("module %s\n\ngo %s\n", request.ModuleName, request.GoVersion)
	if len(requiredImports) > 0 {
		goMod += "\nrequire (\n"
		for _, imp := range requiredImports {
			goMod += fmt.Sprintf("\t%s latest\n", imp)
		}
		goMod += ")\n"
	}
	return goMod
}
