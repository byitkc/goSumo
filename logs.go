package gosumo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type Client struct {
	Endpoint string
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
