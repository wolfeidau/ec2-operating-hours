package operatinghours

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/robfig/cron"
)

var validAction = regexp.MustCompile(`^(?P<action>[a-z]+)=(?P<cron>[0-9*]+ [0-9*]+ [0-9*]+ [0-9*]+ [0-9*]+)$`)

type ActionData struct {
	Action   string
	Cron     string
	Schedule cron.Schedule
}

type InstanceSchedule struct {
	InstanceID string
	Actions    []*ActionData
	Start      *ActionData
	Stop       *ActionData
}

// ParseInstanceTag parse the instance tag and extract actions
func ParseInstanceTag(instanceID, tag string) (*InstanceSchedule, error) {

	tokens := strings.Split(tag, "/")

	//actionData := []*ActionData{}
	instanceSchedule := &InstanceSchedule{InstanceID: instanceID}

	for _, t := range tokens {
		act, err := isValidAction(t)
		if err != nil {
			return nil, err
		}
		//		actionData = append(actionData, act)
		switch act.Action {
		case "start":
			instanceSchedule.Start = act
		case "stop":
			instanceSchedule.Stop = act
		}

	}

	return instanceSchedule, nil
}

func isValidAction(t string) (res *ActionData, err error) {

	match := validAction.FindStringSubmatch(t)
	res = &ActionData{}

	if len(match) == 0 {
		return nil, fmt.Errorf("Invalid action: %s", t)
	}

	for i, name := range validAction.SubexpNames() {

		switch name {
		case "action":
			res.Action = match[i]
		case "cron":
			res.Cron = match[i]
		}

	}

	res.Schedule, err = cron.ParseStandard(res.Cron)

	//spew.Dump(res.Schedule)

	if err != nil {
		return nil, err
	}

	return res, nil
}
