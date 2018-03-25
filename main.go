package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Region      string   `json:"region"`
	Destination []string `json:"destination"`
}

func ParseConfig() (config map[string]Config, err error) {
	config = make(map[string]Config)
	data := make([]byte, 0)
	configURL := os.Getenv("CONFIG_FILE")
	if configURL != "" {
		// Get the data
		resp, err := http.Get(configURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = base64.StdEncoding.DecodeString(os.Getenv("CONFIG"))
		if err != nil {
			return nil, err
		}
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no configuration available")
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ProcessIncomingEvents(request events.S3Event) error {
	config, err := ParseConfig()
	if err != nil {
		return err
	}

	errChan := make(chan error)
	for _, v := range request.Records {
		log.Println("Moving", v.S3.Bucket.Name, v.S3.Object.Key, "To", config[v.S3.Bucket.Name].Destination)
		sess, err := session.NewSession(&aws.Config{Region: aws.String(config[v.S3.Bucket.Name].Region)})
		if err != nil {
			return fmt.Errorf("unable to enstablish aws session for %v", config[v.S3.Bucket.Name])
		}
		for _, v1 := range config[v.S3.Bucket.Name].Destination {
			go copyObjects(s3.New(sess), v.S3.Bucket.Name, v1, v.S3.Object.Key, errChan)
		}
	}

	for _, v := range request.Records {
		for range config[v.S3.Bucket.Name].Destination {
			err = <-errChan
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyObjects(svc *s3.S3, from, to, item string, errChan chan error) {
	_, err := svc.CopyObject(&s3.CopyObjectInput{Bucket: aws.String(to), CopySource: aws.String(from + "/" + item),
		Key: aws.String(item)})
	if err != nil {
		errChan <- fmt.Errorf("unable to copy item %s from bucket %q to bucket %q, %v", item, from, to, err)
		return
	}

	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: aws.String(to), Key: aws.String(item)})
	if err != nil {
		errChan <- fmt.Errorf("error occurred while waiting for item %q to be copied to bucket %q, %v",
			item, to, err)
		return
	}
	errChan <- nil
}

func main() {
	lambda.Start(ProcessIncomingEvents)
}
