package types

import (
	"errors"
	"strings"

	"github.com/hasura/go-graphql-client"
	"github.com/hgiasac/hasura-router/go/types"
)

const (
	ErrCodePermissionDenied = "permission_denied"
	ErrCodeUnsupported      = "unsupported"
)

// NewError create an error instance with code and message
func NewError(code string, message string, extensions map[string]any) types.Error {
	err := types.NewError(code, message)
	if len(extensions) > 0 {
		for k, v := range extensions {
			err.Extensions[k] = v
		}
	}
	return err
}

// ErrUnknown create an error instance with unknown code
func ErrUnknown(err error, extensions map[string]any) types.Error {
	return NewError(types.ErrCodeUnknown, err.Error(), extensions)
}

// ErrUnauthorized create an error instance with unauthorized code
func ErrUnauthorized(err error, extensions map[string]any) types.Error {
	return NewError(types.ErrCodeUnauthorized, err.Error(), extensions)
}

// ErrUnauthorized create an error instance with bad request code
func ErrBadRequest(err error, extensions map[string]any) types.Error {
	return NewError(types.ErrCodeBadRequest, err.Error(), extensions)
}

// ErrDecodeJSON create a bad request error with decode json location
func ErrDecodeJSON(err error, extensions map[string]any) types.Error {
	if extensions == nil {
		extensions = make(map[string]any)
	}
	extensions["location"] = "DecodeJSON"
	return NewError(types.ErrCodeBadRequest, err.Error(), extensions)
}

// ErrValidation create a bad request error with validation location
func ErrValidation(err error, extensions map[string]any) types.Error {
	if extensions == nil {
		extensions = make(map[string]any)
	}
	extensions["location"] = "Validation"
	return NewError(types.ErrCodeBadRequest, err.Error(), extensions)
}

// ErrUnauthorized create an error instance with internal code
func ErrInternal(err error, extensions map[string]any) types.Error {
	return NewError(types.ErrCodeInternal, err.Error(), extensions)
}

// ErrPermissionDenied create an error instance with permission denied code
func ErrPermissionDenied(err error, extensions map[string]any) types.Error {
	message := "permission denied"
	if err != nil {
		message = err.Error()
	}
	return NewError(ErrCodePermissionDenied, message, extensions)
}

// ErrUnsupported create an error instance with unsupported code
func ErrUnsupported(extensions map[string]any) types.Error {
	return NewError(ErrCodeUnsupported, "unsupported", extensions)
}

// Errors the type alias for error slice
type Errors []error

func (errs Errors) Error() string {
	return errs.String()
}

func (errs Errors) String() string {
	var messages []string
	for _, e := range errs {
		messages = append(messages, e.Error())
	}
	return strings.Join(messages, " | ")
}

// ToRouterError tries to convert the error interface to hasura router error
func ToRouterError(err error) error {
	var gqlError graphql.Error
	if errors.As(err, &gqlError) {
		return types.Error{
			Message:    gqlError.Message,
			Extensions: gqlError.Extensions,
		}
	}

	var gqlErrors graphql.Errors
	if errors.As(err, &gqlErrors) && len(gqlErrors) > 0 {
		return types.Error{
			Message:    gqlErrors[0].Message,
			Extensions: gqlErrors[0].Extensions,
		}
	}

	return err
}
