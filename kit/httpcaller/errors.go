package httpcaller

import (
	"errors"
)

// HTTPCallerError is used to pass an error during the request through the
// application with the web specific context.
type HTTPCallerError struct {
	Err      error
	Status   int
	Response []byte
}

// Error implements the error interface for HttpCallerError.
func (hce *HTTPCallerError) Error() string {
	return hce.Err.Error()
}

// NewHTTPCallerError wraps the provided error and its http status, returning an HTTPCallerError.
func NewHTTPCallerError(err error, status int, res []byte) error {
	return &HTTPCallerError{err, status, res}
}

// IsHTTPCallerError is a helper that checks if an error of type HTTPCallerError exists.
func IsHTTPCallerError(err error) bool {
	var hce *HTTPCallerError
	return errors.As(err, &hce)
}

// GetHTTPCallerError returns a copy of the HTTPCallerError pointer.
func GetHTTPCallerError(err error) *HTTPCallerError {
	var hce *HTTPCallerError
	if !errors.As(err, &hce) {
		return nil
	}
	return hce
}
