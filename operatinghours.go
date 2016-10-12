package operatinghours

import (
	"time"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
)

// OperatingHours operating hours
type OperatingHours struct {
	EC2API ec2iface.EC2API
}

type stateChanges struct {
	startInstances []*string
	stopInstances  []*string
}

// NewOperatingHours new operating hours
func NewOperatingHours(sess *session.Session) *OperatingHours {
	svc := ec2.New(sess)
	return &OperatingHours{EC2API: svc}
}

// Check check the operating hours of instances
func (oh *OperatingHours) Check(location string) error {

	loc, err := time.LoadLocation(location)
	if err != nil {
		return err
	}

	stch := &stateChanges{}

	err = oh.EC2API.DescribeInstancesPages(&ec2.DescribeInstancesInput{}, buildDescribeInstanceFunc(loc, stch))
	if err != nil {
		return err
	}

	err = oh.StartInstances(stch.startInstances)
	if err != nil {
		return err
	}

	return oh.StopInstances(stch.stopInstances)
}

// StartInstances start instances in the list provided
func (oh *OperatingHours) StartInstances(startInstances []*string) error {

	log.WithField("len", len(startInstances)).Info("startInstances")

	if len(startInstances) > 0 {
		_, err := oh.EC2API.StartInstances(&ec2.StartInstancesInput{InstanceIds: startInstances})
		if err != nil {
			return err
		}
		for _, instanceID := range startInstances {
			log.WithField("instanceID:", aws.StringValue(instanceID)).Info("start")
		}
	}

	return nil
}

// StopInstances start instances in the list provided
func (oh *OperatingHours) StopInstances(stopInstances []*string) error {

	log.WithField("len", len(stopInstances)).Info("stopInstances")

	if len(stopInstances) > 0 {
		_, err := oh.EC2API.StopInstances(&ec2.StopInstancesInput{InstanceIds: stopInstances})
		if err != nil {
			return err
		}
		for _, instanceID := range stopInstances {
			log.WithField("instanceID:", aws.StringValue(instanceID)).Info("stop")
		}
	}

	return nil
}

func buildDescribeInstanceFunc(loc *time.Location, stch *stateChanges) func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
	return func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {

		for _, instance := range getInstances(page.Reservations) {
			err := processInstance(loc, instance, stch)
			if err != nil {
				return false
			}
		}

		return true
	}
}

func processInstance(loc *time.Location, instance *ec2.Instance, ch *stateChanges) error {
	name := filterTags("Name", instance.Tags)
	operatingHours := filterTags("OperatingHours", instance.Tags)
	state := aws.StringValue(instance.State.Name)
	now := time.Now().In(loc)

	// skip if there is no operating hours tag
	if operatingHours == "" {
		return nil
	}

	sched, err := ParseInstanceTag(aws.StringValue(instance.InstanceId), operatingHours)
	if err != nil {
		return errors.Wrap(err, "error parsing schedule")
	}

	log.WithFields(log.Fields{"name": name, "instanceId": aws.StringValue(instance.InstanceId), "state": state, "start": sched.Start.Cron, "stop": sched.Stop.Cron}).Info("instance")

	performStart, err := TimeToAction(sched.Start.Cron, now, 31*time.Minute)
	if err != nil {
		return errors.Wrap(err, "error actioning schedule")
	}

	log.WithField("performStart", performStart).Info("start check")

	if performStart && state == "stopped" {
		ch.startInstances = append(ch.startInstances, instance.InstanceId)
	}

	log.WithField("len", len(ch.startInstances)).Info("startInstances")

	performStop, err := TimeToAction(sched.Stop.Cron, now, 31*time.Minute)
	if err != nil {
		return errors.Wrap(err, "error actioning schedule")
	}

	log.WithField("performStop", performStop).Info("stop check")

	if performStop && state == "running" {
		ch.stopInstances = append(ch.stopInstances, instance.InstanceId)
	}

	log.WithField("len", len(ch.stopInstances)).Info("stopInstances")

	return nil
}

func getInstances(reservations []*ec2.Reservation) (result []*ec2.Instance) {

	for _, res := range reservations {
		for _, instance := range res.Instances {
			result = append(result, instance)
		}
	}
	return
}

func filterTags(key string, tags []*ec2.Tag) string {

	for _, tag := range tags {
		if aws.StringValue(tag.Key) == key {
			return aws.StringValue(tag.Value)
		}
	}

	return ""
}
