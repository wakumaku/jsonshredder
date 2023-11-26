package forwarder

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// SQSActionResponse SQS Actions and Responses
var SQSActionResponse = map[string]string{
	"GetQueueUrl": responseGetQueueURL,
	"SendMessage": responseSendMessage,
}

const (
	responseGetQueueURL = `
	<GetQueueUrlResponse>
		<GetQueueUrlResult>
		<QueueUrl>https://sqs.us-east-2.amazonaws.com/123456789012/MyQueue</QueueUrl>
		</GetQueueUrlResult>
		<ResponseMetadata>
			<RequestId>470a6f13-2ed9-4181-ad8a-2fdea142988e</RequestId>
			</ResponseMetadata>
			</GetQueueUrlResponse>`
	responseSendMessage = `
			<SendMessageResponse>
		<SendMessageResult>
		<MD5OfMessageBody>5eb63bbbe01eeed093cb22bb8f5acdc3</MD5OfMessageBody>
			<MD5OfMessageAttributes>3ae8f24a165a8cedc005670c81a27295</MD5OfMessageAttributes>
			<MessageId>5fea7756-0ea4-451a-a703-a558b933e274</MessageId>
		</SendMessageResult>
		<ResponseMetadata>
		<RequestId>27daac76-34dd-47df-bd01-1f6e873584a0</RequestId>
		</ResponseMetadata>
	</SendMessageResponse>`
)

// SNSActionResponse SNS Actions and Responses
var SNSActionResponse = map[string]string{
	"Publish": responsePublish,
}

const (
	responsePublish = `
	<PublishResponse xmlns="https://sns.amazonaws.com/doc/2010-03-31/">
		<PublishResult>
			<MessageId>567910cd-659e-55d4-8ccb-5aaf14679dc0</MessageId>
		</PublishResult>
		<ResponseMetadata>
			<RequestId>d74b8436-ae13-5ab4-a9ff-ce54dfea72a0</RequestId>
		</ResponseMetadata>
	</PublishResponse>`
)

// KinesisActionResponse Kinesis Actions and Responses
var KinesisActionResponse = map[string]string{
	"PutRecord": responsePutRecord,
}

const (
	responsePutRecord = `
	{
		"SequenceNumber": "21269319989653637946712965403778482177",
		"ShardId": "shardId-000000000001"
	}`
)

const responseErrorGeneric = `
<ErrorResponse>
<RequestId>42d59b56-7407-4c4a-be0f-4c88daeea257</RequestId>
	<Error>
		<Type>Sender</Type>
		<Code>InvalidParameterValue</Code>
		<Message>
			Value (quename_nonalpha) for parameter QueueName is invalid.
			Must be an alphanumeric String of 1 to 80 in length.
		</Message>
	</Error>
</ErrorResponse>`

func awsMockHandler(actionsResponse map[string]string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		params, err := parseBodyParams(b)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		action := params["Action"]
		if action == "" {
			action = strings.Split(r.Header.Get("X-Amz-Target"), ".")[1]
		}

		response, found := actionsResponse[action]
		if !found {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.URL.Path == "/fail"+action+"/" {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(rw, responseErrorGeneric)
			return
		}

		fmt.Fprint(rw, response)
	}
}

func parseBodyParams(in []byte) (map[string]string, error) {
	r := map[string]string{}

	if err := json.Unmarshal(in, &r); err == nil {
		return r, nil
	}

	values, err := url.ParseQuery(string(in))
	if err != nil {
		return r, err
	}

	for k, v := range values {
		r[k] = v[0]
	}

	return r, nil
}

func startAWSMockServer(kind string) *httptest.Server {
	switch kind {
	case "sns":
		return httptest.NewServer(awsMockHandler(SNSActionResponse))
	case "sqs":
		return httptest.NewServer(awsMockHandler(SQSActionResponse))
	case "kinesis":
		return httptest.NewServer(awsMockHandler(KinesisActionResponse))
	}
	return nil
}
