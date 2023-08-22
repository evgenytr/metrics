// Package logging wraps instantiation of zap.Logger with required configuration.
package logging

import (
	"go.uber.org/zap"
)

// NewDevelopment and NewProduction are supported configs.
const NewDevelopment = "NewDevelopment"
const NewProduction = "NewProduction"

// NewLogger returns zap.Logger with specified config, NewDevelopment by default.
func NewLogger(config string) (logger *zap.Logger, err error) {
	switch config {
	case NewProduction:
		logger, err = zap.NewProduction()
	case NewDevelopment:
		logger, err = zap.NewDevelopment()
	default:
		logger, err = zap.NewDevelopment()
	}
	return
}
