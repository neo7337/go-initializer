package handler

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/neo7337/go-initializer/response"

	"oss.nandlabs.io/golly/rest"
)

type CreateProjectRequest struct {
	ProjectType string `json:"projectType"`
	GoVersion   string `json:"goVersion"`
	Framework   string `json:"framework"`
	ModuleName  string `json:"moduleName"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateProject(ctx rest.ServerContext) {
	var req CreateProjectRequest
	println("[DEBUG] Received CreateProject request")
	if err := json.NewDecoder(ctx.GetRequest().Body).Decode(&req); err != nil {
		println("[ERROR] Invalid request payload:", err.Error())
		response.ResponseJSON(ctx, 400, map[string]string{"error": "Invalid request payload"})
		return
	}

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	println("[DEBUG] Created zip writer")

	// Example: Add a README.md file to the zip (customize as needed)
	f, err := zipWriter.Create("README.md")
	if err != nil {
		println("[ERROR] Failed to create README.md in zip:", err.Error())
		response.ResponseJSON(ctx, 500, map[string]string{"error": "Failed to create zip file"})
		return
	}
	println("[DEBUG] Created README.md in zip")
	_, err = f.Write([]byte("# " + req.Name + "\n" + req.Description))
	if err != nil {
		println("[ERROR] Failed to write to README.md in zip:", err.Error())
		response.ResponseJSON(ctx, 500, map[string]string{"error": "Failed to write to zip file"})
		return
	}
	println("[DEBUG] Wrote to README.md in zip")

	// Add more files as needed based on req

	if err := zipWriter.Close(); err != nil {
		println("[ERROR] Failed to finalize zip file:", err.Error())
		response.ResponseJSON(ctx, 500, map[string]string{"error": "Failed to finalize zip file"})
		return
	}
	println("[DEBUG] Finalized zip file, size:", buf.Len())

	ctx.SetHeader("Content-Type", "application/zip")
	ctx.SetHeader("Content-Disposition", "attachment; filename=project.zip")
	ctx.SetStatusCode(http.StatusOK)
	println("[DEBUG] Sending zip file response")
	ctx.Write(buf.Bytes(), "application/zip")
}
