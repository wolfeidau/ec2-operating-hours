package operatinghours

import (
	"fmt"
	"time"
)

const shortForm = "2006-Jan-02 15:04"

func parseTestDate(dt string) time.Time {
	t, _ := time.ParseInLocation(shortForm, dt, time.Local)
	return t
}

func parseDate(dt string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, dt)
	fmt.Println(err)
	return t
}
