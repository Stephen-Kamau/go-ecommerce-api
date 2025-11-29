package utils

import "fmt"

// ---------------------
// Resource / Entity Errors
// ---------------------

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID '%s' not found", e.Resource, e.ID)
}

// AlreadyExistsError represents an attempt to create a resource that already exists
type AlreadyExistsError struct {
	Resource string
	ID       string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s with ID '%s' already exists", e.Resource, e.ID)
}

// ---------------------
// Validation Errors
// ---------------------

// ValidationError represents a failure in input validation
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s': %s", e.Field, e.Message)
}

// ---------------------
// Authentication / Authorization Errors
// ---------------------

// AuthenticationError represents login/auth failure
type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string {
	if e.Message == "" {
		return "authentication failed"
	}
	return e.Message
}

// AuthorizationError represents access denied
type AuthorizationError struct {
	Action string
}

func (e *AuthorizationError) Error() string {
	return fmt.Sprintf("not authorized to perform action: %s", e.Action)
}

// ---------------------
// Database / Persistence Errors
// ---------------------

// DatabaseError wraps generic database errors
type DatabaseError struct {
	Query string
	Err   error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error on query '%s': %v", e.Query, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// ---------------------
// External API / Service Errors
// ---------------------

// ExternalServiceError represents an error calling an external service
type ExternalServiceError struct {
	Service string
	Err     error
}

func (e *ExternalServiceError) Error() string {
	return fmt.Sprintf("external service '%s' error: %v", e.Service, e.Err)
}

func (e *ExternalServiceError) Unwrap() error {
	return e.Err
}

// ---------------------
// Internal / Unexpected Errors
// ---------------------

// InternalError represents an unexpected internal error
type InternalError struct {
	Message string
	Err     error
}

func (e *InternalError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("internal error: %s - %v", e.Message, e.Err)
	}
	return fmt.Sprintf("internal error: %s", e.Message)
}

func (e *InternalError) Unwrap() error {
	return e.Err
}
