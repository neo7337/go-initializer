package main

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeZipBuffer(t *testing.T, entries map[string]string) *bytes.Buffer {
	t.Helper()
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for path, content := range entries {
		w, err := zw.Create(path)
		require.NoError(t, err)
		_, err = w.Write([]byte(content))
		require.NoError(t, err)
	}
	require.NoError(t, zw.Close())
	return buf
}

func makeRunNewCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "new"}
	cmd.Flags().String("description", "", "")
	cmd.Flags().StringArray("addon", nil, "")
	cmd.Flags().Bool("docker", false, "")
	_ = cmd.Flags().Parse([]string{"--description=d", "--docker=false", "--addon=cache=redis"})
	return cmd
}

func TestParseAddons_ValidAndRepeatCategories(t *testing.T) {
	in := []string{"cache=redis", "database=gorm", "cache=memcached"}
	out, err := parseAddons(in)
	require.NoError(t, err)

	require.Contains(t, out, "cache")
	require.Contains(t, out, "database")
	assert.ElementsMatch(t, []string{"redis", "memcached"}, out["cache"])
	assert.Equal(t, []string{"gorm"}, out["database"])
}

func TestParseAddons_InvalidFormat(t *testing.T) {
	for _, tc := range []string{"cache", "=redis", "cache=", ""} {
		t.Run(tc, func(t *testing.T) {
			_, err := parseAddons([]string{tc})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid --addon value")
		})
	}
}

func TestStripFirstComponent(t *testing.T) {
	assert.Equal(t, "cmd/main.go", stripFirstComponent("myapp/cmd/main.go"))
	assert.Equal(t, "", stripFirstComponent("myapp"))
	assert.Equal(t, "b/c", stripFirstComponent("a/b/c"))
}

func TestBuildSelectOpts_SortedAscending(t *testing.T) {
	opts := buildSelectOpts(map[string]bool{"z": true, "a": true, "m": true})
	require.Len(t, opts, 3)
	assert.Contains(t, opts[0].String(), "a")
	assert.Contains(t, opts[1].String(), "m")
	assert.Contains(t, opts[2].String(), "z")
}

func TestBuildSelectOptsDesc_SortedDescending(t *testing.T) {
	opts := buildSelectOptsDesc(map[string]bool{"1.24.6": true, "1.25.0": true, "1.23.12": true})
	require.Len(t, opts, 3)
	assert.Contains(t, opts[0].String(), "1.25.0")
	assert.Contains(t, opts[1].String(), "1.24.6")
	assert.Contains(t, opts[2].String(), "1.23.12")
}

func TestBuildLabelSelectOpts_UsesLabelAndSortedKeys(t *testing.T) {
	opts := buildLabelSelectOpts(map[string]string{
		"microservice": "Microservice",
		"api-server":   "API Server",
	})
	require.Len(t, opts, 2)
	assert.Contains(t, opts[0].String(), "API Server")
	assert.Contains(t, opts[1].String(), "Microservice")
}

func TestExtractZip_ExtractsIntoOutputRoot(t *testing.T) {
	tmp := t.TempDir()
	buf := makeZipBuffer(t, map[string]string{
		"myapp/README.md":         "# myapp",
		"myapp/cmd/myapp/main.go": "package main",
	})

	err := extractZip(buf, tmp)
	require.NoError(t, err)

	readmePath := filepath.Join(tmp, "README.md")
	mainPath := filepath.Join(tmp, "cmd", "myapp", "main.go")
	_, err = os.Stat(readmePath)
	require.NoError(t, err)
	_, err = os.Stat(mainPath)
	require.NoError(t, err)
}

func TestExtractZip_RejectsZipSlip(t *testing.T) {
	tmp := t.TempDir()
	buf := makeZipBuffer(t, map[string]string{
		"myapp/../../evil.txt": "nope",
	})

	err := extractZip(buf, tmp)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "would escape output directory")
}

func TestExtractZip_InvalidZip(t *testing.T) {
	err := extractZip(bytes.NewBufferString("not a zip"), t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read zip")
}

func TestPrintSuccess_UsesBuildForCliApp(t *testing.T) {
	cmd := &cobra.Command{}
	out := new(bytes.Buffer)
	cmd.SetOut(out)

	printSuccess(cmd, "mycli", "cli-app", "/tmp/mycli")
	body := out.String()
	assert.Contains(t, body, "Project created at /tmp/mycli")
	assert.Contains(t, body, "cd mycli")
	assert.Contains(t, body, "go mod tidy")
	assert.Contains(t, body, "make build")
	assert.NotContains(t, body, "make run")
}

func TestPrintSuccess_UsesRunForServiceTypes(t *testing.T) {
	cmd := &cobra.Command{}
	out := new(bytes.Buffer)
	cmd.SetOut(out)

	printSuccess(cmd, "svc", "microservice", "/tmp/svc")
	body := out.String()
	assert.Contains(t, body, "make run")
	assert.NotContains(t, body, "make build")
}

func TestRunNew_UnknownProjectType(t *testing.T) {
	cmd := makeRunNewCmd()
	opts := &newOpts{
		name:        "x",
		module:      "github.com/acme/x",
		goVersion:   "1.24.6",
		projectType: "not-a-real-type",
		framework:   "gin",
		output:      t.TempDir(),
	}

	err := runNew(cmd, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown project type")
}

func TestRunNew_SimpleProjectSuccess(t *testing.T) {
	outDir := t.TempDir()
	cmd := makeRunNewCmd()
	out := new(bytes.Buffer)
	cmd.SetOut(out)

	opts := &newOpts{
		name:        "tp",
		module:      "github.com/acme/tp",
		description: "test project",
		goVersion:   "1.24.6",
		projectType: "simple-project",
		framework:   "golly",
		output:      outDir,
	}

	err := runNew(cmd, opts)
	require.NoError(t, err)

	// runNew extracts generated content directly into outDir (without top-level folder).
	_, err = os.Stat(filepath.Join(outDir, "README.md"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(outDir, "go.mod"))
	require.NoError(t, err)

	body := out.String()
	assert.Contains(t, body, "Generating simple-project project")
	assert.Contains(t, body, "Project created at")
}
