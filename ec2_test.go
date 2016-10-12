package operatinghours

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type instanceData struct {
	instanceID     string
	operatingHours string
	execTime       time.Time
	performAction  bool
}

var ()

func TestParseStart(t *testing.T) {

	instancesStart := []instanceData{
		{"i-123123123", "start=0 6 * * */stop=0 20 * * *", parseTestDate("2013-Feb-03 05:45"), true},
		{"i-123123124", "start=0 7 * * */stop=0 21 * * *", parseTestDate("2013-Feb-03 06:45"), true},
		{"i-123123125", "start=0 8 * * */stop=0 22 * * *", parseTestDate("2013-Feb-03 06:45"), false},
		{"i-123123125", "start=0 8 * * */stop=0 22 * * *", parseTestDate("2013-Feb-03 10:45"), false},
	}

	for _, instance := range instancesStart {

		instanceSchedule, err := ParseInstanceTag(instance.instanceID, instance.operatingHours)
		assert.Nil(t, err)

		performStart, err := TimeToAction(instanceSchedule.Start.Cron, instance.execTime, 31*time.Minute)
		assert.Nil(t, err)

		//fmt.Printf("performAction = %v, performStop = %v\n", instance.performAction, performStart)
		assert.Equal(t, instance.performAction, performStart)
	}
}

func TestParseStop(t *testing.T) {
	instancesStop := []instanceData{
		{"i-123123123", "start=* 6 * * */stop=* 20 * * *", parseTestDate("2013-Feb-03 19:45"), true},
		{"i-123123124", "start=* 7 * * */stop=* 21 * * *", parseTestDate("2013-Feb-03 20:45"), true},
		{"i-123123125", "start=* 8 * * */stop=* 22 * * *", parseTestDate("2013-Feb-03 23:00"), false},
		{"i-123123125", "start=* 8 * * */stop=* 22 * * *", parseTestDate("2013-Feb-04 01:45"), false},
	}

	for _, instance := range instancesStop {

		instanceSchedule, err := ParseInstanceTag(instance.instanceID, instance.operatingHours)
		assert.Nil(t, err)

		performStop, err := TimeToAction(instanceSchedule.Stop.Cron, instance.execTime, 31*time.Minute)
		assert.Nil(t, err)
		//fmt.Printf("performAction = %v, performStop = %v\n", instance.performAction, performStop)
		assert.Equal(t, instance.performAction, performStop)

	}

}
