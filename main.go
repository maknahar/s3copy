package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func ParseConfig() (config map[string][]string, err error) {
	data := make([]byte, 0)
	configURL := os.Getenv("CONFIG_FILE")
	if configURL == "" {
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

	config = make(map[string][]string)

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

	var wg sync.WaitGroup
	for _, v := range request.Records {
		region, err := s3manager.GetBucketRegion(nil, session.Must(session.NewSession()), v.S3.Bucket.Name, "ap-southeast-1")
		if err != nil {
			return err
		}

		sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
		for _, v1 := range config[v.S3.Bucket.Name] {
			wg.Add(1)
			go copyObjects(&wg, s3.New(sess), v.S3.Bucket.Name, v1, v.S3.Object.Key)
		}
	}
	wg.Wait()
	return nil
}

func copyObjects(wg *sync.WaitGroup, svc *s3.S3, from, to, item string) error {
	defer wg.Done()
	_, err := svc.CopyObject(&s3.CopyObjectInput{Bucket: aws.String(to), CopySource: aws.String(from + "/" + item),
		Key: aws.String(item)})
	if err != nil {
		return fmt.Errorf("unable to copy item %s from bucket %q to bucket %q, %v", item, from, to, err)
	}

	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: aws.String(to), Key: aws.String(item)})
	if err != nil {
		return fmt.Errorf("error occurred while waiting for item %q to be copied to bucket %q, %v",
			item, to, err)
	}

	return nil
}

func main() {
	lambda.Start(ProcessIncomingEvents)
}
