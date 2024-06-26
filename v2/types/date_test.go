package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestDate(t *testing.T) {

	for _, ip := range []struct {
		Input   json.RawMessage
		IsError bool
		Output  *Date
	}{
		{
			[]byte("\"2020-01-01\""),
			false,
			&Date{2020, 1, 1},
		}, {
			[]byte("\"2020-02-29\""),
			false,
			&Date{2020, 2, 29},
		}, {
			[]byte("\"2020-02-30\""),
			true,
			nil,
		}, {
			[]byte("\"\""),
			true,
			nil,
		},
	} {
		t.Run(string(ip.Input), func(t *testing.T) {
			var result Date
			err := json.Unmarshal(ip.Input, &result)
			if ip.IsError {
				assert.Error(t, err, fmt.Sprintf("invalid date `%s`", strings.Trim(string(ip.Input), "\"")))
			} else {
				assert.NilError(t, err)
				assert.Equal(t, *ip.Output, result)
			}
		})
	}

	t.Run("test date object", func(t *testing.T) {
		var dateObject struct {
			Date *Date `json:"date"`
		}

		err := json.Unmarshal([]byte(`{"date": null}`), &dateObject)
		assert.NilError(t, err)
		assert.Assert(t, dateObject.Date == nil)
	})
}
