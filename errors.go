package telegram

import (
	"fmt"
	"strings"
)

var errUnauthorized = fmt.Errorf("unauthorized")

// IsUnauthorizedError checks if error is unauthorized
func IsUnauthorizedError(err error) bool {
	return err == errUnauthorized
}

var errForbidden = fmt.Errorf("forbidden")

// IsForbiddenError checks if error is forbidden
func IsForbiddenError(err error) bool {
	return err == errForbidden
}

// IsAPIError checks if error is ApiError
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

// IsRequiredError checks if error is RequiredError
func IsRequiredError(err error) bool {
	_, ok := err.(*RequiredError)
	return ok
}

// IsValidationError checks if error is ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// APIError contains error information from response
type APIError struct {
	Description string `json:"description"`
	// ErrorCode contents are subject to change in the future.
	ErrorCode int `json:"error_code"`
}

// Error returns string representation for ApiError
func (e *APIError) Error() string {
	return fmt.Sprintf("apiError: %s", e.Description)
}

// RequiredError tells if fields are required but were not filled
type RequiredError struct {
	Fields []string
}

// Error returns string representation for RequiredError
func (e *RequiredError) Error() string {
	return fmt.Sprintf("%s required", strings.Join(e.Fields, " or "))
}

// NewRequiredError creates RequireError
func NewRequiredError(fields ...string) *RequiredError {
	return &RequiredError{Fields: fields}
}

// NewValidationError creates ValidationError
func NewValidationError(field string, description string) *ValidationError {
	return &ValidationError{
		Field:       field,
		Description: description,
	}
}

// ValidationError tells if field has wrong value
type ValidationError struct {
	// Field name
	Field       string `json:"field"`
	Description string `json:"description"`
}

// Error returns string representation for ValidationError
func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		"field %s is invalid: %s",
		e.Field,
		e.Description)
}
