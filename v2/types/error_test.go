package types

import (
	"errors"
	"testing"

	"github.com/hasura/go-graphql-client"
	"github.com/hgiasac/hasura-router/go/types"
	"gotest.tools/v3/assert"
)

func TestError_ToRouterError(t *testing.T) {
	assert.DeepEqual(t, types.Error{
		Message: "test",
		Extensions: map[string]any{
			"code": "unknown",
		},
	}, ToRouterError(graphql.Error{
		Message: "test",
		Extensions: map[string]any{
			"code": "unknown",
		},
	}, nil))

	assert.DeepEqual(t, types.Error{
		Message: "test",
		Extensions: map[string]any{
			"code": "unknown",
			"foo":  "bar",
		},
	}, ToRouterError(graphql.Errors{
		{
			Message: "test",
			Extensions: map[string]any{
				"code": "unknown",
			},
		},
	}, map[string]any{
		"foo": "bar",
	}))

	assert.ErrorContains(t, ToRouterError(errors.New("test"), nil), "test")
	assert.ErrorContains(t, ToRouterError(errors.New("test"), map[string]any{"foo": "bar"}), "unknown: test; extensions: map[code:unknown foo:bar]")
}
