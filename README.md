# copys3
A lambda function to copy objects from n buckets
 to m buckets in same region.

## Config
lambda Function expect the list if source and 
destination buckets in a JSON file.


### Config File Format  
```
{
  "source_bucket": [
    "destination bucket1",
    "destination bucket2"
  ]
}
```

### Env Var

CONFIG_FILE : URL of configuration file