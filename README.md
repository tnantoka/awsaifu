# AWSSaifu

## Requirements

- https://github.com/aws/aws-sdk-go/

## Example

```
$ go get github.com/tnantoka/awsaifu

# ./hosts
default

# ~/.aws/credentials
[default]
aws_access_key_id = 
aws_secret_access_key = 

# ./.api_token
API_TOKEN

$ awsaifu -r ROOM_ID
[info][title]AWSの課金額[/title]
default: 160ドル
[/info]
```

