package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresArray(t *testing.T) {

	arr, err := DecodePostgresArray("{a,b,c}")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, arr, []string{"a", "b", "c"}, "DecodePostgresArray")
	assert.Equal(t, EncodePostgresArray(arr), "{a,b,c}", "DecodePostgresArray")
}
