# copys3
A lambda function to copy objects from n buckets
 to m buckets in same region.

## Config
This Lambda Function expect the list of source and 
destination buckets in a JSON file.


### Config File Format  
```
{
  "s3copyinput": {
    "region": "us-east-1",
    "destination": [
      "s3copyoutput1",
      "s3copyoutput2"
    ]
  }
}
```

### Env Var
This Lambda function can be configured in two ways. 
Either give a public url of config file via CONFIG_FILE or
provide base64 encoded value of configuration via CONFIG.

If both are provided CONFIG_FILE value will take precedence.
 

CONFIG_FILE : URL of configuration file

CONFIG : Base64 Encoded string of content of CONFIG_FILE

### NOTE:
Make sure Lambda function have required access to source and
destination bucket. 