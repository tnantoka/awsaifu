package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/service/ec2"
	"flag"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	. "github.com/tnantoka/chatsworth"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func main() {
	var p = flag.String("p", "./profiles", "AWS Profiles")
	var k = flag.String("k", "./.api_token", "Chatwork API Token")
	var r = flag.String("r", "", "ChatWork Room ID")
	flag.Parse()

	cw := Chatsworth{
		RoomID:   *r,
		APIToken: loadToken(*k),
	}
	cw.PostMessage(buildMessage(*p))
}

func loadToken(file string) string {
	token, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(token)
}

func buildMessage(file string) string {
	profiles, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	var validProfiles []string
	for _, profile := range strings.Split(string(profiles), "\n") {
		if len(profile) > 0 {
			validProfiles = append(validProfiles, profile)
		}
	}

	messageChan := fetchCharges(validProfiles)

	message := "[info][title]AWSの課金額[/title]"
	for i := 0; i < len(validProfiles); i++ {
		m := <-messageChan
		fmt.Print(m)
		message += m
	}
	message += "[/info]"

	return message
}

func fetchCharges(profiles []string) <-chan string {
	messageChan := make(chan string)

	for _, profile := range profiles {
		go func(profile string) {
			config := aws.Config{Region: "us-east-1"}
			config.Credentials = credentials.NewSharedCredentials("", profile)
			message := profile + ": " + fetchCharge(config) + "ドル\n"
			messageChan <- message
		}(profile)
	}

	return messageChan
}

func fetchCharge(config aws.Config) string {
	dimension := cloudwatch.Dimension{
		Name:  aws.String("Currency"),
		Value: aws.String("USD"),
	}

	svc := cloudwatch.New(&config)
	input := cloudwatch.GetMetricStatisticsInput{
		Dimensions: []*cloudwatch.Dimension{&dimension},
		StartTime:  aws.Time(time.Now().Add(-24 * time.Hour)),
		EndTime:    aws.Time(time.Now()),
		MetricName: aws.String("EstimatedCharges"),
		Namespace:  aws.String("AWS/Billing"),
		Period:     aws.Long(60),
		Statistics: []*string{aws.String("Maximum")},
		//Unit: "",
	}
	output, err := svc.GetMetricStatistics(&input)

	if err != nil {
		log.Fatal(err)
	}

	var dp = output.Datapoints[0]
	return fmt.Sprint(*dp.Maximum)
}
