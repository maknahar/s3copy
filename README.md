# copys3
A lambda function to copy objects from n buckets
 to m buckets in same region.
 
 
## Deployment
- Download the deployment guide from 
https://github.com/maknahar/copys3/releases
- Upload the deployment file to lambda
- Set the configuration

or

- Checkout the latest release code `go get github.com/maknahar/copys3`
- Build the deployment file `GOOS=linux go build -o main main.go && zip deployment.zip main`
- Upload the deployment file to lambda
- Set the configuration 

## Config
This Lambda Function expect the list of source and 
destination buckets in a JSON file.


### Config File Format  

Config File is a map of input s3 bucket and detail to copy 
changes to destination.

Function can take event from SQS as well if you have the 
source bucket put all event to a SQS.

Function can take event from SQS where SQS is subscribe to
a SNS and source bucket publish a message to SNS on every 
change.

To Support SQS processing add SQS URL `sqs` and
SQS region `sqsRegion` in the config file. Otherwise,
Leave them empty `""` or remove the keys.

```
{
  "s3copyinput": {
    "region": "us-east-1",
    "sqs": "url of SQS where all S3 change events are stored",
    "sqsRegion": "us-east-1"
    "destinations": [
      "s3copyoutput1",
      "s3copyoutput2"
    ]
  }
}
```

### Env Var
This Lambda function can be configured in two ways. 
- Either give a public url of config file via CONFIG_FILE
- or provide base64 encoded value of configuration via 
CONFIG.

If both are provided CONFIG_FILE value will take precedence.

CONFIG_FILE : URL of configuration file

CONFIG : Base64 Encoded string of content of CONFIG_FILE

## NOTE:
Make sure Lambda function have required access to source 
and destination bucket. 

## Contribution
All contributions are welcome. Either via PR or Issue.