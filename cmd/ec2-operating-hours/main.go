package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/aws/aws-sdk-go/aws/session"
	operatinghours "github.com/wolfeidau/ec2-operating-hours"
)

func main() {
	log.SetHandler(cli.Default)

	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("error building aws session: %s", err)
	}

	oh := operatinghours.NewOperatingHours(sess)

	err = oh.Check("Australia/Melbourne")
	if err != nil {
		log.Fatalf("error processing operating hours: %v", err)
	}

}
