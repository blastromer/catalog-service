package apierror

import "fmt"

// Kind classifies an application error for HTTP mapping.
type Kind string

const (
	KindValidation Kind = "validation"
	KindNotFound   Kind = "not_found"
	KindConflict   Kind = "conflict"
	KindInternal   Kind = "internal"
)

// Error is a typed application error returned by services.
type Error struct {
	Kind    Kind
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func Validation(message string) *Error {
	return &Error{Kind: KindValidation, Message: message}
}

func NotFound(message string) *Error {
	return &Error{Kind: KindNotFound, Message: message}
}

func Internal(message string, err error) *Error {
	return &Error{Kind: KindInternal, Message: message, Err: err}
}

// As extracts an *Error from err when present.
func As(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	if apiErr, ok := err.(*Error); ok {
		return apiErr, true
	}
	return nil, false
}
