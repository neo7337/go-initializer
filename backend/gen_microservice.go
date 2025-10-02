package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"log"
)

func GenerateMicroservice(request CreateProjectRequest) (*bytes.Buffer, error) {
	// Implementation for generating a microservice project based on input
	// generate a zip file in memory

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Create a folder with the project name and write files inside it
	folderName := request.Name
	if folderName == "" {
		folderName = "project"
	}

	// Example: add a README.md file inside the folder
	readmeContent := fmt.Sprintf("# %s\n\n%s", request.Name, request.Description)
	readmeFile, err := zipWriter.Create(fmt.Sprintf("%s/README.md", folderName))
	if err != nil {
		log.Printf("[ERROR] Failed to create file in zip: %v", err)
		err = errors.New("failed to create file in zip")
		return nil, err
	}
	_, err = readmeFile.Write([]byte(readmeContent))
	if err != nil {
		log.Printf("[ERROR] Failed to write to zip: %v", err)
		err = errors.New("failed to write to zip file")
		return nil, err
	}

	// Add more files as needed based on request
	err = zipWriter.Close()
	if err != nil {
		log.Printf("[ERROR] Failed to close zip writer: %v", err)
		err = errors.New("failed to finalize zip file")
		return nil, err
	}

	return buf, nil
}
