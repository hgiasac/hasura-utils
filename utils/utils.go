package utils

import (
	"fmt"
	"strings"
)

// EncodePostgresArray encode array string to postgres array
func EncodePostgresArray(input []string) string {
	return fmt.Sprintf("{%s}", strings.Join(input, ","))
}

// DecodePostgresArray decode postgres array string
func DecodePostgresArray(input string) ([]string, error) {
	if input == "{}" {
		return []string{}, nil
	}

	if len(input) < 3 || input[0] != '{' || input[len(input)-1] != '}' {
		return nil, fmt.Errorf("invalid postgres array: %s", input)
	}

	return strings.Split(input[1:len(input)-1], ","), nil
}

// EncodePostgresArrayStringer the generic function to encode postgres array from Stringer implementations
func EncodePostgresArrayStringer[V fmt.Stringer](inputs []V) string {
	length := len(inputs)
	if length == 0 {
		return "{}"
	}

	sInputs := make([]string, length)
	for i, u := range inputs {
		sInputs[i] = u.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(sInputs, ","))
}
