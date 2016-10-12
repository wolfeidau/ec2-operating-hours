package operatinghours

import (
	"time"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

// DayDuration duration for 24 hours
const DayDuration = 24 * time.Hour

// InTimeSpan test if the check date is in the period
func InTimeSpan(start time.Time, duration time.Duration, check time.Time) bool {
	//fmt.Printf("start = %v, end = %v, check = %v\n", start, start.Add(duration), check)
	log.WithFields(log.Fields{"start": start, "end": start.Add(duration), "check": check}).Info("InTimeSpan")
	return check.After(start) && check.Before(start.Add(duration))
}

// TimeToAction check time till the action is performed
func TimeToAction(sched string, now time.Time, duration time.Duration) (bool, error) {
	cron, err := cron.ParseStandard(sched)
	if err != nil {
		return false, errors.Wrap(err, "Unable to parse schedule")
	}
	next := cron.Next(now)

	return InTimeSpan(now, duration, next), nil
}
