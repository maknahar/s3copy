package main

import (
	"testing"
)

func TestProcessIncomingEvents(t *testing.T) {
	//events.S3Event{
	//	Records: []events.S3EventRecord{{
	//		S3: events.S3Entity{
	//			Bucket: events.S3Bucket{
	//				Name: "s3copyinput",
	//			},
	//			Object: events.S3Object{
	//				Key: "config.json",
	//			},
	//		},
	//	}},
	//}
	m := make(map[string]interface{})
	err := ProcessIncomingEvents(m)
	if err != nil {
		t.Fatal("Got Error ", err)
	}
}
