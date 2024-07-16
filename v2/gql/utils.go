package gql

import (
	"context"
	"fmt"

	"github.com/hasura/go-graphql-client"
)

type contextHeaderKey string

const headerKey contextHeaderKey = "x-headers"

const (
	XHasuraUserID                    = "x-hasura-user-id"
	XHasuraAdminSecret               = "x-hasura-admin-secret"
	XHasuraRole                      = "x-hasura-role"
	XRequestId                       = "x-request-id"
	XHasuraUseBackendOnlyPermissions = "x-hasura-use-backend-only-permissions"

	RoleAdmin string = "admin"
)

// getHeadersFromContext get request headers from the context
func getHeadersFromContext(ctx context.Context) map[string]string {
	h := ctx.Value(headerKey)
	var headers map[string]string
	if h == nil {
		headers = map[string]string{}
	} else {
		headers = h.(map[string]string)
	}
	return headers
}

func setHeaders(ctx context.Context, hs map[string]string) context.Context {
	h := ctx.Value(headerKey)
	var headers map[string]string
	if h == nil {
		headers = map[string]string{}
	} else {
		headers = h.(map[string]string)
	}
	for k, v := range hs {
		headers[k] = v
	}
	return context.WithValue(ctx, headerKey, headers)
}

func getOperationNameFromOptions(options []graphql.Option) string {
	for _, opt := range options {
		if opt.Type() == "operation_name" {
			if stringer, ok := opt.(fmt.Stringer); ok {
				return stringer.String()
			}
			break
		}
	}
	return ""
}
