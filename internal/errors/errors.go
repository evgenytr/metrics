package errors

import (
	"time"

	"github.com/evgenytr/metrics.git/internal/utils"
)

var RepeatedAttemptsIntervals = [3]*time.Duration{
	utils.GetTimeInterval(1),
	utils.GetTimeInterval(3),
	utils.GetTimeInterval(5),
}
