package types

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Date represents a date object which is compatible with Postgres date type
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// Date returns the time.Time instance from date
func (dt Date) Date() time.Time {
	return time.Date(dt.Year, dt.Month, dt.Day, 0, 0, 0, 0, time.UTC)
}

// MarshalXML implements the xml Marshaler interface
func (dt Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(dt.String(), start)
}

// UnmarshalXML implements the xml Unmarshaler interface
func (dt *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	if s != "" && strings.ToLower(s) != "null" {
		r, err := ParseDate(s)
		if err != nil {
			return err
		}
		*dt = *r
	}

	return nil
}

// UnmarshalJSON implements the json Unmarshaler interface
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	result, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = *result
	return nil
}

// MarshalJSON implements the json Marshaler interface
func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// String implements the Stringer interface
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// MustParseDate parses date from string, panic if error
func MustParseDate(input string) *Date {
	r, err := ParseDate(input)
	if err != nil {
		panic(err)
	}
	return r
}

// ParseDate parses date from string
func ParseDate(input string) (*Date, error) {
	parts := strings.Split(input, "-")
	invalidError := fmt.Errorf("invalid date `%s`", input)

	if len(parts) != 3 {
		return nil, invalidError
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil || year < 0 || year > 9999 {
		return nil, invalidError
	}

	if year < 0 {
		return nil, invalidError
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil || month <= 0 || month > 12 {
		return nil, invalidError
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil || day <= 0 ||
		((month == 1 || month == 3 || month == 5 || month == 7 || month == 8 || month == 10 || month == 12) && day > 31) ||
		((month == 4 || month == 6 || month == 9 || month == 11) && day > 30) ||
		(month == 2 && year%4 > 0 && day > 28) ||
		(month == 2 && year%4 == 0 && day > 29) {
		return nil, invalidError
	}

	return &Date{
		Year:  year,
		Month: time.Month(month),
		Day:   day,
	}, nil
}
