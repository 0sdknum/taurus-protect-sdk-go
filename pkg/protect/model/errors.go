package model

import (
	"fmt"
	"time"
)

// IntegrityError indicates a cryptographic verification failure.
// This is a security-critical error that should never be retried.
type IntegrityError struct {
	Message string
	Err     error
}

func (e *IntegrityError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("integrity error: %s", e.Message)
	}
	return "integrity error"
}

func (e *IntegrityError) Unwrap() error {
	return e.Err
}

// WhitelistError indicates a whitelist verification failure.
type WhitelistError struct {
	Message string
	Err     error
}

func (e *WhitelistError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("whitelist error: %s", e.Message)
	}
	return "whitelist error"
}

func (e *WhitelistError) Unwrap() error {
	return e.Err
}

// APIError represents an error response from the Taurus-PROTECT API.
type APIError struct {
	Message     string
	Code        int
	StatusCode  int
	Description string
	ErrorCode   string
	// ResponseBody is the raw API error response body, when available.
	ResponseBody string
	Err          error
	RetryAfter   time.Duration
}

func (e *APIError) Error() string {
	code := e.HTTPStatusCode()
	if e.Description != "" && e.Message != "" {
		return fmt.Sprintf("%s: %s (code=%d)", e.Description, e.Message, code)
	}
	if e.Message != "" {
		return fmt.Sprintf("%s (code=%d)", e.Message, code)
	}
	if e.Description != "" {
		return fmt.Sprintf("%s (code=%d)", e.Description, code)
	}
	return fmt.Sprintf("API error (code=%d)", code)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// As exposes the canonical APIError through specialized errors that embed it.
func (e *APIError) As(target any) bool {
	apiError, ok := target.(**APIError)
	if !ok {
		return false
	}
	*apiError = e
	return true
}

// HTTPStatusCode returns the HTTP status code for this error.
func (e *APIError) HTTPStatusCode() int {
	if e.Code != 0 {
		return e.Code
	}
	return e.StatusCode
}

// IsRetryable returns true if the error is retryable (429 or 5xx).
func (e *APIError) IsRetryable() bool {
	code := e.HTTPStatusCode()
	return code == 429 || (code >= 500 && code < 600)
}

// IsClientError returns true for 4xx errors.
func (e *APIError) IsClientError() bool {
	code := e.HTTPStatusCode()
	return code >= 400 && code < 500
}

// IsServerError returns true for 5xx errors.
func (e *APIError) IsServerError() bool {
	code := e.HTTPStatusCode()
	return code >= 500 && code < 600
}

// SuggestedRetryDelay returns the server-provided delay or a default retry delay.
func (e *APIError) SuggestedRetryDelay() time.Duration {
	code := e.HTTPStatusCode()
	if code == 429 {
		if e.RetryAfter > 0 {
			return e.RetryAfter
		}
		return time.Second
	}
	if code >= 500 {
		return 5 * time.Second
	}
	return 0
}

// Is matches API errors by HTTP status code category.
func (e *APIError) Is(target error) bool {
	targetError, ok := target.(*APIError)
	if !ok {
		return false
	}
	code := e.HTTPStatusCode()
	switch targetError.HTTPStatusCode() {
	case 400:
		return code == 400
	case 401:
		return code == 401
	case 403:
		return code == 403
	case 404:
		return code == 404
	case 429:
		return code == 429
	case 500:
		return code >= 500
	default:
		return code == targetError.HTTPStatusCode()
	}
}

// ValidationError represents a 400 Bad Request error.
type ValidationError struct {
	*APIError
}

// AuthenticationError represents a 401 Unauthorized error.
type AuthenticationError struct {
	*APIError
}

// AuthorizationError represents a 403 Forbidden error.
type AuthorizationError struct {
	*APIError
}

// NotFoundError represents a 404 Not Found error.
type NotFoundError struct {
	*APIError
}

// RateLimitError represents a 429 Too Many Requests error.
type RateLimitError struct {
	*APIError
}

// ServerError represents a 5xx server error.
type ServerError struct {
	*APIError
}

// ConfigurationError represents a client configuration error.
type ConfigurationError struct {
	Message string
	Err     error
}

func (e *ConfigurationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("configuration error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("configuration error: %s", e.Message)
}

func (e *ConfigurationError) Unwrap() error {
	return e.Err
}

// RequestMetadataError represents an error related to request metadata.
type RequestMetadataError struct {
	Message string
	Err     error
}

func (e *RequestMetadataError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("request metadata error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("request metadata error: %s", e.Message)
}

func (e *RequestMetadataError) Unwrap() error {
	return e.Err
}

// NewAPIError creates the appropriate typed error based on the HTTP status code.
func NewAPIError(statusCode int, errorCode string, message string, err error) error {
	base := &APIError{
		Code:       statusCode,
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
		Err:        err,
	}
	switch {
	case statusCode == 400:
		return &ValidationError{APIError: base}
	case statusCode == 401:
		return &AuthenticationError{APIError: base}
	case statusCode == 403:
		return &AuthorizationError{APIError: base}
	case statusCode == 404:
		return &NotFoundError{APIError: base}
	case statusCode == 429:
		return &RateLimitError{APIError: base}
	case statusCode >= 500:
		return &ServerError{APIError: base}
	default:
		return base
	}
}
