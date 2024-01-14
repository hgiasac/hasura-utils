package types

import (
	"errors"
	"testing"

	"github.com/hasura/go-graphql-client"
	"github.com/hgiasac/hasura-router/go/types"
	"github.com/stretchr/testify/assert"
)

func TestError_ToRouterError(t *testing.T) {
	assert.Equal(t, types.Error{
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

	assert.Equal(t, types.Error{
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

	assert.EqualError(t, ToRouterError(errors.New("test"), nil), "test")
	assert.EqualError(t, ToRouterError(errors.New("test"), map[string]any{"foo": "bar"}), "unknown: test; extensions: map[code:unknown foo:bar]")
}
