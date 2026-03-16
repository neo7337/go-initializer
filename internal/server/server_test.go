package server

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/neo7337/go-initializer/internal/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_IsInitializedAndValidatesStruct(t *testing.T) {
	require.NotNil(t, validate)

	invalid := generator.CreateProjectRequest{}
	err := validate.Struct(invalid)
	require.Error(t, err)

	valid := generator.CreateProjectRequest{
		ProjectType: "simple-project",
		GoVersion:   "1.24.6",
		Framework:   "golly",
		ModuleName:  "github.com/acme/x",
		Name:        "x",
	}
	err = validate.Struct(valid)
	require.NoError(t, err)
}

func TestStart_BuildsExpectedServerConfig(t *testing.T) {
	origListen := listenAndServe
	origFatal := logFatalf
	defer func() {
		listenAndServe = origListen
		logFatalf = origFatal
	}()

	var captured *http.Server
	listenAndServe = func(s *http.Server) error {
		captured = s
		return nil
	}
	logFatalf = func(string, ...interface{}) {
		t.Fatalf("logFatalf should not be called on successful startup")
	}

	Start()

	require.NotNil(t, captured)
	assert.Equal(t, ":8182", captured.Addr)
	assert.Equal(t, 10*time.Second, captured.ReadTimeout)
	assert.Equal(t, 10*time.Second, captured.WriteTimeout)
	assert.Equal(t, 1<<20, captured.MaxHeaderBytes)
	require.NotNil(t, captured.Handler)
}

func TestStart_ListenErrorCallsFatalf(t *testing.T) {
	origListen := listenAndServe
	origFatal := logFatalf
	defer func() {
		listenAndServe = origListen
		logFatalf = origFatal
	}()

	listenErr := errors.New("boom")
	listenAndServe = func(*http.Server) error { return listenErr }

	called := false
	var format string
	logFatalf = func(f string, args ...interface{}) {
		called = true
		format = f
		if len(args) == 1 {
			assert.Equal(t, listenErr, args[0])
		}
	}

	Start()

	assert.True(t, called, "expected logFatalf to be called")
	assert.Contains(t, format, "Server failed to start")
}
