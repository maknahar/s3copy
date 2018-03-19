package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func ProcessIncomingEvents(request events.S3Event) error {
	// Get the data
	resp, err := http.Get(os.Getenv("CONFIG_FILE"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	config := make(map[string][]string)

	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, v := range request.Records {
		sess := session.Must(session.NewSession())
		region, err := s3manager.GetBucketRegion(nil, sess, v.S3.Bucket.Name, "ap-southeast-1")
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
				fmt.Fprintf(os.Stderr, "unable to find bucket %s's region not found\n", v.S3.Bucket.Name)
			}
			return err
		}

		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		svc := s3.New(sess)

		for _, v1 := range config[v.S3.Bucket.Name] {
			wg.Add(1)
			go copyObjects(&wg, svc, v.S3.Bucket.Name, v1, v.S3.Object.Key)
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

	// Wait to see if the item got copied
	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: aws.String(to), Key: aws.String(item)})
	if err != nil {
		return fmt.Errorf("error occurred while waiting for item %q to be copied to bucket %q, %v",
			item, to, err)
	}

	log.Printf("Item %q successfully copied from bucket %q to bucket %q\n", item, from, to)
	return nil
}

func main() {
	lambda.Start(ProcessIncomingEvents)
	//log.Println(ProcessIncomingEvents())
}
