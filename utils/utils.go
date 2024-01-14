package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	Digits        = "0123456789"
	Alphabets     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphaDigits   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var src = rand.NewSource(time.Now().UnixNano())

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

// GenRandomString generate random string with fixed length
func GenRandomString(n int, allowedCharacters ...string) string {
	allowedChars := AlphaDigits
	if len(allowedCharacters) > 0 {
		allowedChars = allowedCharacters[0]
	}
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(allowedChars) {
			sb.WriteByte(allowedChars[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
