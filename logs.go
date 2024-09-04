package gosumo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type Client struct {
	Endpoint string
}

// NewClient creates and returns a new client using the provided endpointURL.
// It will check to ensure that the URL is valid and will return an error if it
// is not.
func NewClient(endpointURL string) (Client, error) {
	if _, err := url.Parse(endpointURL); err != nil {
		return Client{}, ErrBuildingClient{
			Message: fmt.Sprintf("unable to build client using the URL '%s'", endpointURL),
		}
	}
	return Client{Endpoint: endpointURL}, nil
}

// PostLogs will post the logs provided as a slice of logs. All logs structs
// must include Metadata for JSON encoding.
// It will return an error if there are problems parsing or posting the logs to
// the Sumo Logic Endpoint.
func PostLogs[T any](c Client, logs []T) error {
	sLogs, err := getJSONString(logs)
	if err != nil {
		return ErrParsingLogs{
			Message: fmt.Sprintf("error parsing logs: %v", err),
		}
	}
	if err := PostLogsString(c, sLogs); err != nil {
		return ErrPostingLogs{
			Message: err.Error(),
		}
	}
	return nil
}

// PostLogsString will post the logs provided as a string (newline separated) to
// the provided Sumo Logic Client Endpoint.
// The provided logs can be in any format, and should be delimited with a \n
// (newline character).
func PostLogsString(c Client, logs string) error {
	logReader := strings.NewReader(logs)
	req, err := http.NewRequest("POST", c.Endpoint, logReader)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrPostingLogs{
			Message: fmt.Sprintf("unexpected status code when posting logs, expected: %d, got: %d", http.StatusOK, resp.StatusCode),
		}
	}
	return nil
}

// getJSONString takes a slice of structs that include JSON metadata. It returns
// a string with all JSON objects as a string containing all logs delimited by a
// newline character (\n).
func getJSONString[T any](s []T) (string, error) {
	var sLogs []string
	for _, v := range s {
		if !hasJSONMetadata(v) {
			return "", ErrParsingLogs{
				Message: fmt.Sprint("object is missing json metadata"),
			}
		}
		bLog, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		sLogs = append(sLogs, string(bLog))
	}
	return strings.Join(sLogs, "\n"), nil
}

// hasJSONMetadata takes a struct and checks to confirm that all values inside
// of the struct have JSON metadata for Marshalling before posting to Sumo Logic.
func hasJSONMetadata(a any) bool {
	val := reflect.ValueOf(a)
	if val.Kind() != reflect.Struct {
		return false
	}
	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			return false
		}
	}
	return true
}
