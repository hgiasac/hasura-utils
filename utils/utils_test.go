package utils

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestPostgresArray(t *testing.T) {

	arr, err := DecodePostgresArray("{a,b,c}")
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, arr, []string{"a", "b", "c"})
	assert.DeepEqual(t, EncodePostgresArray(arr), "{a,b,c}")
}
