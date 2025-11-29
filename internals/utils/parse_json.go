package utils

// to parse JSON from http request body
import (
	"encoding/json"
	"io"
)

func ParseJSON(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(dest)
	if err != nil {
		return err
	}
	return nil
}
