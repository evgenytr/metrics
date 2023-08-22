// Package errors contains error handling helpers.
package errors

import (
	"time"

	"github.com/evgenytr/metrics.git/internal/utils"
)

// RepeatedAttemptsIntervals array contains predefined intervals for repeat attempts
var RepeatedAttemptsIntervals = [3]*time.Duration{
	utils.GetTimeInterval(1),
	utils.GetTimeInterval(3),
	utils.GetTimeInterval(5),
}
