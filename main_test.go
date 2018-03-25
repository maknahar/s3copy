package main

import (
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

func TestProcessIncomingEvents(t *testing.T) {
	err := ProcessIncomingEvents(events.S3Event{
		Records: []events.S3EventRecord{{
			S3: events.S3Entity{
				Bucket: events.S3Bucket{
					Name: "s3copyinput",
				},
				Object: events.S3Object{
					Key: "config.json",
				},
			},
		}},
	})
	if err != nil {
		t.Fatal("Got Error ", err)
	}
}
