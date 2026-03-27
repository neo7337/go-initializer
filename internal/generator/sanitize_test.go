package generator

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── ValidateProjectName ──────────────────────────────────────────────────────

func TestValidateProjectName_Valid(t *testing.T) {
	valid := []string{
		"myapp",
		"my-app",
		"my_app",
		"App123",
		"a",
		strings.Repeat("a", 64), // exactly 64 chars (1 leading + 63 more)
	}
	for _, name := range valid {
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, ValidateProjectName(name))
		})
	}
}

func TestValidateProjectName_RejectsForwardSlash(t *testing.T) {
	err := ValidateProjectName("my/app")
	require.Error(t, err)
	var ve *ErrValidation
	require.True(t, errors.As(err, &ve))
	assert.Equal(t, "name", ve.Field)
}

func TestValidateProjectName_RejectsBackslash(t *testing.T) {
	err := ValidateProjectName("my\\app")
	require.Error(t, err)
	var ve *ErrValidation
	require.True(t, errors.As(err, &ve))
	assert.Equal(t, "name", ve.Field)
}

func TestValidateProjectName_RejectsDotDot(t *testing.T) {
	err := ValidateProjectName("..")
	require.Error(t, err)
	var ve *ErrValidation
	require.True(t, errors.As(err, &ve))
	assert.Equal(t, "name", ve.Field)
}

func TestValidateProjectName_RejectsNullByte(t *testing.T) {
	err := ValidateProjectName("app\x00name")
	require.Error(t, err)
	var ve *ErrValidation
	require.True(t, errors.As(err, &ve))
	assert.Equal(t, "name", ve.Field)
}

func TestValidateProjectName_RejectsLeadingDigit(t *testing.T) {
	err := ValidateProjectName("1app")
	require.Error(t, err)
	var ve *ErrValidation
	require.True(t, errors.As(err, &ve))
	assert.Equal(t, "name", ve.Field)
}

func TestValidateProjectName_RejectsTooLong(t *testing.T) {
	// 65 chars: 1 valid leading letter + 64 more
	err := ValidateProjectName("a" + strings.Repeat("b", 64))
	require.Error(t, err)
	var ve *ErrValidation
	require.True(t, errors.As(err, &ve))
	assert.Equal(t, "name", ve.Field)
}

func TestValidateProjectName_RejectsLeadingDot(t *testing.T) {
	err := ValidateProjectName(".hidden")
	require.Error(t, err)
	var ve *ErrValidation
	require.True(t, errors.As(err, &ve))
}

// ─── ValidateModuleName ───────────────────────────────────────────────────────

func TestValidateModuleName_Valid(t *testing.T) {
	valid := []string{
		"github.com/acme/myapp",
		"example.com/foo/bar-baz",
		"mymodule",
		"github.com/acme/myapp/v2",
	}
	for _, m := range valid {
		t.Run(m, func(t *testing.T) {
			assert.NoError(t, ValidateModuleName(m))
		})
	}
}

func TestValidateModuleName_RejectsBackslash(t *testing.T) {
	err := ValidateModuleName("github.com\\acme\\myapp")
	require.Error(t, err)
	var ve *ErrValidation
	require.True(t, errors.As(err, &ve))
	assert.Equal(t, "moduleName", ve.Field)
}

func TestValidateModuleName_RejectsNullByte(t *testing.T) {
	err := ValidateModuleName("github.com/acme/\x00app")
	require.Error(t, err)
	var ve *ErrValidation
	require.True(t, errors.As(err, &ve))
	assert.Equal(t, "moduleName", ve.Field)
}

func TestValidateModuleName_RejectsDotDotSegment(t *testing.T) {
	err := ValidateModuleName("github.com/../secret")
	require.Error(t, err)
	var ve *ErrValidation
	require.True(t, errors.As(err, &ve))
	assert.Equal(t, "moduleName", ve.Field)
}

func TestValidateModuleName_RejectsTooLong(t *testing.T) {
	// 257 chars total
	err := ValidateModuleName("a" + strings.Repeat("b", 256))
	require.Error(t, err)
	var ve *ErrValidation
	require.True(t, errors.As(err, &ve))
	assert.Equal(t, "moduleName", ve.Field)
}

// ─── ErrValidation.Error() message format ────────────────────────────────────

func TestErrValidation_ErrorMessage(t *testing.T) {
	e := &ErrValidation{Field: "name", Message: "bad input"}
	assert.Equal(t, "invalid name: bad input", e.Error())
}
