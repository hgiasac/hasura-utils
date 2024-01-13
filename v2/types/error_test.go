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
	}))

	assert.Equal(t, types.Error{
		Message: "test",
		Extensions: map[string]any{
			"code": "unknown",
		},
	}, ToRouterError(graphql.Errors{
		{
			Message: "test",
			Extensions: map[string]any{
				"code": "unknown",
			},
		},
	}))

	assert.EqualError(t, ToRouterError(errors.New("test")), "test")
}
