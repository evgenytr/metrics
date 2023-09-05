// Package utils contains helper functions.
package utils

import "time"

func GetTimeInterval(seconds float64) time.Duration {
	value := time.Duration(seconds) * time.Second
	return value
}
