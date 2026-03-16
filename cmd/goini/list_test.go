package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	require.NoError(t, w.Close())
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	require.NoError(t, r.Close())

	return buf.String()
}

func TestNewListCmd_HasExpectedSubcommands(t *testing.T) {
	cmd := newListCmd()
	names := make([]string, 0, len(cmd.Commands()))
	for _, sub := range cmd.Commands() {
		names = append(names, sub.Name())
	}

	assert.Contains(t, names, "types")
	assert.Contains(t, names, "frameworks")
	assert.Contains(t, names, "addons")
}

func TestListTypesCommand_PrintsSupportedTypes(t *testing.T) {
	cmd := newListTypesCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	out := captureStdout(t, func() {
		cmd.SetArgs([]string{})
		err := cmd.Execute()
		require.NoError(t, err)
	})

	assert.Contains(t, out, "TYPE")
	assert.Contains(t, out, "LABEL")
	assert.Contains(t, out, "api-server")
	assert.Contains(t, out, "cli-app")
	assert.Contains(t, out, "microservice")
	assert.Contains(t, out, "simple-project")
	assert.Contains(t, out, "API Server")
	assert.Contains(t, out, "CLI Application")
}

func TestListFrameworksCommand_RequiresTypeFlag(t *testing.T) {
	cmd := newListFrameworksCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--type is required")
}

func TestListFrameworksCommand_UnknownType(t *testing.T) {
	cmd := newListFrameworksCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.SetArgs([]string{"--type", "unknown"})
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown project type")
}

func TestListFrameworksCommand_PrintsFrameworksForType(t *testing.T) {
	cmd := newListFrameworksCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	out := captureStdout(t, func() {
		cmd.SetArgs([]string{"--type", "microservice"})
		err := cmd.Execute()
		require.NoError(t, err)
	})

	assert.Contains(t, out, "FRAMEWORK")
	assert.Contains(t, out, "gin")
	assert.Contains(t, out, "echo")
	assert.Contains(t, out, "fiber")
	assert.Contains(t, out, "gokit")
	assert.Contains(t, out, "golly")
}

func TestListFrameworksCommand_OutputIsSorted(t *testing.T) {
	cmd := newListFrameworksCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	out := captureStdout(t, func() {
		cmd.SetArgs([]string{"--type", "cli-app"})
		err := cmd.Execute()
		require.NoError(t, err)
	})

	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.GreaterOrEqual(t, len(lines), 6)
	frameworkLines := lines[2:]

	joined := strings.Join(frameworkLines, "\n")
	cobraIdx := strings.Index(joined, "cobra")
	gollyIdx := strings.Index(joined, "golly")
	kingpinIdx := strings.Index(joined, "kingpin")
	urfaveIdx := strings.Index(joined, "urfave")

	require.NotEqual(t, -1, cobraIdx)
	require.NotEqual(t, -1, gollyIdx)
	require.NotEqual(t, -1, kingpinIdx)
	require.NotEqual(t, -1, urfaveIdx)
	assert.Less(t, cobraIdx, gollyIdx)
	assert.Less(t, gollyIdx, kingpinIdx)
	assert.Less(t, kingpinIdx, urfaveIdx)
}

func TestListAddonsCommand_PrintsAllCategoriesAndAddons(t *testing.T) {
	cmd := newListAddonsCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	out := captureStdout(t, func() {
		cmd.SetArgs([]string{})
		err := cmd.Execute()
		require.NoError(t, err)
	})

	assert.Contains(t, out, "CATEGORY")
	assert.Contains(t, out, "ADDON")
	assert.Contains(t, out, "cache")
	assert.Contains(t, out, "redis")
	assert.Contains(t, out, "memcached")
	assert.Contains(t, out, "database")
	assert.Contains(t, out, "gorm")
	assert.Contains(t, out, "ent")
	assert.Contains(t, out, "other")
	assert.Contains(t, out, "zap")
	assert.Contains(t, out, "logrus")
	assert.Contains(t, out, "cobra")
}

func TestListAddonsCommand_OutputStartsWithAlphabeticalCategory(t *testing.T) {
	cmd := newListAddonsCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	out := captureStdout(t, func() {
		cmd.SetArgs([]string{})
		err := cmd.Execute()
		require.NoError(t, err)
	})

	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.GreaterOrEqual(t, len(lines), 3)
	firstDataLine := lines[2]
	assert.Contains(t, firstDataLine, "cache")
}
