package monitor

import (
	"errors"
	"net/http"
	"time"
)

type (
	StatusFunc func(statusCode int, body map[string]interface{}) error
	Monitor    struct {
		Name          string
		Interval      time.Duration
		URL           string
		statusFunc    StatusFunc
		Retries       int
		RetryInterval time.Duration
		Timeout       time.Duration
		onFailure     func()
	}
)

var DefaultStatusFunc StatusFunc = func(statusCode int, body map[string]interface{}) error {
	if statusCode != http.StatusOK {
		return errors.New("status code was not 200")
	}

	if body["status"] != "UP" {
		return errors.New("status was not 'UP'")
	}
	return nil
}
