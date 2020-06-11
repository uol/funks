package funks

import (
	"encoding/json"
	"errors"
	"time"
)

// Duration - a duration wrapper type to add the method below
type Duration struct {
	time.Duration
}

// UnmarshalText - used by the toml parser to proper parse duration values
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// MarshalJSON - transforms the value to the encoded json format
func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON - transforms the json format to the duration format
func (d *Duration) UnmarshalJSON(b []byte) error {

	var v interface{}
	var err error

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {

	case float64:

		d.Duration = time.Duration(value)

		return nil

	case string:

		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}

		return nil

	default:
		return errors.New("invalid duration")
	}
}
