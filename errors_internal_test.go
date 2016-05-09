package telegram

import (
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func TestIsUnauthorizedError(t *testing.T) {
	assert.True(t, IsUnauthorizedError(errUnauthorized))
	assert.False(t, IsUnauthorizedError(errForbidden))
}

func TestIsForbiddenError(t *testing.T) {
	assert.True(t, IsForbiddenError(errForbidden))
	assert.False(t, IsForbiddenError(errUnauthorized))
}

func TestIsAPIError(t *testing.T) {
	assert.True(t, IsAPIError(&APIError{}))
	assert.False(t, IsAPIError(errUnauthorized))
}

func TestIsRequiredError(t *testing.T) {
	assert.True(t, IsRequiredError(&RequiredError{}))
	assert.False(t, IsRequiredError(errUnauthorized))
}

func TestIsValidationError(t *testing.T) {
	assert.True(t, IsValidationError(&ValidationError{}))
	assert.False(t, IsValidationError(errUnauthorized))
}
