# Local resources for testing

# SQS
echo "Creating SQS queues ..."
awslocal sqs create-queue --queue-name queue
awslocal sqs create-queue --queue-name queue2

# SNS
echo "Creating SNS topics ..."
awslocal sns create-topic --name topic

# Kinesis Datastream
echo "Creating Kinesis Datastreams ..."
awslocal kinesis create-stream --stream-name stream --shard-count 1

# DynamoDB table


echo "done."