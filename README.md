# JSON Shredder (Oroku Saki)

Transforms JSON objects to another objects with the same values.

Have you ever ... (jo mai mai ...)

- ... received big payloads where you just need a couple of fields?
- ... wanted build new endpoints without waiting for the caller re-implementation?
- ... received complex nested payloads and you just need a key/value flatten object?

**jsonshredder can help you,**
**receive the payload you deserve!**

- Define transformations to create **new payloads** from old ones
- Use it as **transformation** service
- Use it as **forwarder** service (transform and forward)

## Example case

If you want to transform this:

```json
{
    "payload": {
        "user": {
            "data":{
                "user_id": 12345,
                "city":    "NY",
                "geo":     {"lat": 1, "long": 1},
                "user":    "john.doe",
                "address": {
                    "street": "foo avenue",
                    "number": 42
                },
                "email":    "john.doe@bar.tld"
            }
        }
    }
}
```

To this:

```json
{
    "username":      "john.doe",
    "email_address": "john.doe@bar.tld"
}
```

You can use this transformation:

```yaml
# ...
    mappings:
      - path: payload.user.data.user
        path_out: username
      - path: payload.user.data.email
        path_out: email_address
```

How?

### Transformation

1. Create a config file named `myconfig.yml`:

```yaml
---
port: 8080

transformations:
  - name: myfirsttransformation
    mappings:
      - path: payload.user.data.user
        path_out: username
      - path: payload.user.data.email
        path_out: email_address
```

2. Run jsonshredder passing this config:

```shell
$ docker run --rm -p 8080:8080 -v $(pwd)/myconfig.yml:/config.yml wakumaku/jsonshredder:latest
...
```

3. Test it:

```shell
$ curl -X POST -d '{"payload":{"user":{"data":{"user_id":12345,"city":"NY","geo":{"lat":1,"long":1},"user":"john.doe","address":{"street":"foo avenue","number":42},"email":"john.doe@bar.tld"}}}}' http://localhost:8080/myfirsttransformation

{"email_address":"john.doe@bar.tld","username":"john.doe"}
```

### Transformation and forwarding

Now imagine that you want it's to forward the resulting transformation to an AWS SQS Queue.

1. Add an SQS forwarder configuration in your config file like this:

```yaml
---
port: 8080

transformations:
  - name: myfirsttransformation
    mappings:
      - path: payload.user.data.user
        path_out: username
      - path: payload.user.data.email
        path_out: email_address

forwarders:
  - name: mysqs
    kind: sqs
    params:
      aws_profile: mydelegatedrole
      aws_region: us-east-1
      aws_resource_name: "queue"
```

2. Run the jsonshredder:

```shell
$ docker run --rm -p 8080:8080 \
  -v ~/.aws:/.aws:ro \
  -e AWS_SHARED_CREDENTIALS_FILE=/.aws/credentials \
  -v $(pwd)/config.dev.yaml:/config.yml \
  wakumaku/jsonshredder:latest
...
```

NOTE: here we are sharing our aws credentials because we are using a profile.

If you prefer, you can specify the `aws_access_key_id` and `aws_secret_access_key` params.


3. Test it:

```shell
$ curl -X POST -d '{"payload":{"user":{"data":{"user_id":12345,"city":"NY","geo":{"lat":1,"long":1},"user":"john.doe","address":{"street":"foo avenue","number":42},"email":"john.doe@bar.tld"}}}}' http://localhost:8080/myfirsttransformation/mysqs

{"email_address":"john.doe@bar.tld","username":"john.doe"}
```

Take a look to the URL: `http://localhost:8080/myfirsttransformation/mysqs`

Note that the first segment of the path is the Transformation name and the second one the Forwarder name.

This way you can combine Transformations and Forwardings as you wish!

4. Check for new items in your queue!

```shell
$ aws sqs get-message --queue-url http://...
``` 

## Configuration file

### Global

- `port`: HTTP Server port number
- `loglevel`: debug, info, warn or err

### Transformations

Each transformation MUST have a unique name. Those names will be used in the HTTP Endpoints, please use `[A-Za-z0-9_-]` chars.

```yaml
  - name: client_payment     # name that will be used in the URL to identify the transformation
    operation: extract       # add | extract, default: extract. Add appends new fields to the original structure.
    mappings:
      - path: user.data.age  # jmespath expression: https://jmespath.org/
        path_out: age        # path to create with the value (or object) found in path
        type_out: string     # string | int | float, default: original value type. Can force conversions. Adds quotes when is string and try to convert to int/float.
        default_null: 0      # default value when the path doesn't exist or the value is null
```

### Forwarders

There is a limited of forwarders implemented, currently: http, sns, sqs, kinesis. They can be used to simplify the transformation result delivery.

AWS Family forwarders:

```yaml
  - name: myforwarder # unique name to identify the forwarder
    kind: sns         # sns, sqs, kinesis
    params:
      aws_endpoint: http://localstack:4566 # allows point to an AWS Mock service
      aws_access_key_id: foo               # AWS KeyID
      aws_secret_access_key: bar           # AWS SecretAccessKey
      aws_profile: profile                 # AWS Profile name
      aws_resource_arn: aws:::resource/... # When kind is: SNS
      aws_resource_name: resource_name     # When kind is: SQS, KINESIS
      aws_region: us-east-1                # AWS Region
```

Other forwarders:

```yaml
  - name: myforwarder
    kind: http
    params:
      http_endpoint: http://....com/post # endpoint destination
      http_header_auth: Bearer mytoken   # (optional) Authorization header value
      http_status_ok: 200                # expected status code from the destination
```

## Full Configuration example

```yaml
### Global
port: 8080
loglevel: debug

### Transformations
transformations:
  - name: myfirsttransformation
    mappings:
      - path: payload.user.data.user
        path_out: username
      - path: payload.user.data.email
        path_out: email_address

  - name: transform1
    mappings:
      - path: user.data.name
        path_out: username
      - path: user.data.email
        path_out: email

  - name: transform2
    operation: add
    mappings:
      - path: username
        path_out: username2
      - path: email
        path_out: email2
        default_null: "empty@email.com"

  - name: transform3
    operation: add
    mappings:
      - path: username2
        path_out: username3
      - path: email2
        path_out: email3

  - name: transform4
    operation: add
    mappings:
      - path: username3
        path_out: username4
      - path: email3
        path_out: email4

### Forwarders
forwarders:

  - name: sendtosns
    kind: sns
    params:
      aws_endpoint: http://localstack:4566
      aws_access_key_id: foo
      aws_secret_access_key: bar
      aws_region: us-east-1
      aws_resource_arn: "arn:aws:sns:us-east-1:000000000000:topic"

  - name: sendtosqs
    kind: sqs
    params:
      aws_endpoint: http://localstack:4566
      aws_access_key_id: foo
      aws_secret_access_key: bar
      aws_region: us-east-1
      aws_resource_name: "queue"

  - name: sendtosqs2
    kind: sqs
    params:
      aws_endpoint: http://localstack:4566
      aws_access_key_id: foo
      aws_secret_access_key: bar
      aws_region: us-east-1
      aws_resource_name: "queue2"

  - name: sendtokinesis
    kind: kinesis
    params:
      aws_endpoint: http://localstack:4566
      aws_access_key_id: foo
      aws_secret_access_key: bar
      aws_region: us-east-1
      aws_resource_name: "stream"

  - name: sendtohttp
    kind: http
    params:
      http_endpoint: http://localhost:8080/transform2/sendtohttp2
      http_header_auth: Bearer mytoken
      http_status_ok: 200

  - name: sendtohttp2
    kind: http
    params:
      http_endpoint: http://localhost:8080/transform3/sendtohttp3
      http_header_auth: Bearer mytoken
      http_status_ok: 200

  - name: sendtohttp3
    kind: http
    params:
      http_endpoint: http://localhost:8080/transform4/sendtosqs
      http_header_auth: Bearer mytoken
      http_status_ok: 200
```

Try this config chaining HTTP Forwarders!

```shell
$ curl -X POST -H 'Content-type:application/json' -d '{"user":{"data":{"name":"john"}}}'  http://localhost:8080/transform1/sendtohttp
```

This is the result you'll find in the queue:

```json
{
    "email": null,
    "email2": "empty@email.com",
    "email3": "empty@email.com",
    "email4": "empty@email.com",
    "username": "john",
    "username2": "john",
    "username3": "john",
    "username4": "john"
}
```