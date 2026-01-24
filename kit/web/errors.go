package web

import "errors"

// RequestError is used to pass an error during the request through the
// application with the web specific context.
type RequestError struct {
	Err    error
	Status int
}

// Error implements the error interface for RequestError.
func (re *RequestError) Error() string {
	return re.Err.Error()
}

// NewRequestError wraps the provided error and its http status, returning a RequestError.
func NewRequestError(err error, status int) error {
	return &RequestError{err, status}
}

// IsRequestError is a helper that checks if an error of type RequestError exists.
func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}

// GetRequestError returns a copy of the RequestError pointer.
func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}

type BadGatewayError struct {
	Err     error
	Message string
}

func (e BadGatewayError) Error() string {
	return e.Err.Error()
}

func IsBadGatewayError(err error) bool {
	var e *BadGatewayError
	return errors.As(err, &e)
}

func NewBadGatewayError(err error, message string) error {
	return BadGatewayError{Err: err, Message: message}
}

func GetBadGatewayError(err error) *BadGatewayError {
	var e *BadGatewayError
	if !errors.As(err, &e) {
		return nil
	}
	return e
}

type ConflictError struct {
	Err           error
	CustomMessage string
}

// ConflictError wraps the provided error and its http status, returning a ConflictError.
func NewConflictError(err error, msg string) error {
	return &ConflictError{err, msg}
}

// Error implements the error interface for ConflictError.
func (re *ConflictError) Error() string {
	return re.Err.Error()
}

// IsConflictError is a helper that checks if an error of type ConflictError exists.
func IsConflictError(err error) bool {
	var re *ConflictError
	return errors.As(err, &re)
}

// GetConflictError returns a copy of the ConflictError pointer.
func GetConflictError(err error) *ConflictError {
	var re *ConflictError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}
