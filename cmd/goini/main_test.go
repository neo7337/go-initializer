package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureStdoutMain(t *testing.T, fn func()) string {
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

func TestNewRootCmd_HasExpectedMetadataAndSubcommands(t *testing.T) {
	cmd := newRootCmd()

	assert.Equal(t, "goini", cmd.Use)
	assert.Equal(t, "Go project scaffolding tool", cmd.Short)

	names := make([]string, 0, len(cmd.Commands()))
	for _, sub := range cmd.Commands() {
		names = append(names, sub.Name())
	}
	assert.Contains(t, names, "version")
	assert.Contains(t, names, "completion")
	assert.Contains(t, names, "list")
	assert.Contains(t, names, "new")
}

func TestNewVersionCmd_PrintsVersion(t *testing.T) {
	cmd := newVersionCmd()
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	orig := version
	version = "v9.9.9-test"
	defer func() { version = orig }()

	out := captureStdoutMain(t, func() {
		cmd.SetArgs([]string{})
		err := cmd.Execute()
		require.NoError(t, err)
	})

	assert.Contains(t, out, "goini v9.9.9-test")
}

func TestNewCompletionCmd_RequiresOneValidArg(t *testing.T) {
	root := newRootCmd()
	cmd := newCompletionCmd(root)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 1 arg")
}

func TestNewCompletionCmd_RejectsInvalidArg(t *testing.T) {
	root := newRootCmd()
	cmd := newCompletionCmd(root)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	cmd.SetArgs([]string{"invalid-shell"})
	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid argument")
}

func TestNewCompletionCmd_GeneratesBashCompletion(t *testing.T) {
	root := newRootCmd()
	cmd := newCompletionCmd(root)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	out := captureStdoutMain(t, func() {
		cmd.SetArgs([]string{"bash"})
		err := cmd.Execute()
		require.NoError(t, err)
	})

	assert.Contains(t, out, "__start_goini")
}

func TestNewCompletionCmd_GeneratesZshCompletion(t *testing.T) {
	root := newRootCmd()
	cmd := newCompletionCmd(root)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	out := captureStdoutMain(t, func() {
		cmd.SetArgs([]string{"zsh"})
		err := cmd.Execute()
		require.NoError(t, err)
	})

	assert.Contains(t, out, "#compdef goini")
}

func TestNewCompletionCmd_GeneratesFishCompletion(t *testing.T) {
	root := newRootCmd()
	cmd := newCompletionCmd(root)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	out := captureStdoutMain(t, func() {
		cmd.SetArgs([]string{"fish"})
		err := cmd.Execute()
		require.NoError(t, err)
	})

	assert.Contains(t, out, "complete -c goini")
}

func TestNewCompletionCmd_GeneratesPowerShellCompletion(t *testing.T) {
	root := newRootCmd()
	cmd := newCompletionCmd(root)
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true

	out := captureStdoutMain(t, func() {
		cmd.SetArgs([]string{"powershell"})
		err := cmd.Execute()
		require.NoError(t, err)
	})

	assert.Contains(t, out, "Register-ArgumentCompleter")
	assert.Contains(t, out, "goini")
}
