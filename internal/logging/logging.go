package logging

import (
	"go.uber.org/zap"
)

const NewDevelopment = "NewDevelopment"
const NewProduction = "NewProduction"

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
