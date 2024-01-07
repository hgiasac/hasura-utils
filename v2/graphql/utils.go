package graphql

import "context"

type contextHeaderKey string

const headerKey contextHeaderKey = "x-headers"

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
