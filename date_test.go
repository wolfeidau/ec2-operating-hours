package operatinghours

import (
	"testing"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetHandler(cli.Default)

}

func TestWithin(t *testing.T) {
	execs := []struct {
		start    time.Time
		duration time.Duration
		check    time.Time
		running  bool
	}{
		{parseTestDate("2013-Feb-03 05:45"), 31 * time.Minute, parseTestDate("2013-Feb-03 05:50"), true},
		{parseTestDate("2014-Feb-03 05:45"), 31 * time.Minute, parseTestDate("2014-Feb-03 22:00"), false},
		{
			start:    parseDate("2016-10-06T18:32:35.016678586+11:00"),
			duration: 31 * time.Minute,
			check:    parseDate("2016-10-06T19:00:00+11:00"),
			running:  true,
		},
	}

	for _, exec := range execs {
		assert.Equal(t, exec.running, InTimeSpan(exec.start, exec.duration, exec.check))
	}
}
