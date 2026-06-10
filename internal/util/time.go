package util

import (
	"errors"
	"time"
)

const DateLayout = "2006-01-02"

func NowString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func ValidateDate(value string) error {
	if value == "" {
		return errors.New("date is required")
	}
	_, err := time.Parse(DateLayout, value)
	return err
}
