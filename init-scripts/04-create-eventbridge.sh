#!/bin/bash
set -euo pipefail

echo "Setting up EventBridge..."

EVENT_BUS_NAME="default"

echo "Enabling EventBridge notifications on S3 bucket..."
awslocal s3api put-bucket-notification-configuration \
    --bucket "${AWS_S3_BUCKET}" \
    --notification-configuration '{
        "EventBridgeConfiguration": {}
    }'

echo "Creating EventBridge rule..."
awslocal events put-rule \
    --name S3UploadRule \
    --event-bus-name "${EVENT_BUS_NAME}" \
    --event-pattern '{
        "source": ["aws.s3"],
        "detail-type": ["Object Created"]
    }' \
    --state ENABLED

echo "Getting SQS Queue ARN..."
QUEUE_ARN=$(awslocal sqs get-queue-attributes \
    --queue-url "http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/${SQS_NOTIFICATIONS_QUEUE}" \
    --attribute-names QueueArn \
    --query 'Attributes.QueueArn' \
    --output text)

echo "SQS Queue ARN: ${QUEUE_ARN}"

echo "Setting SQS as target for EventBridge rule..."
awslocal events put-targets \
    --rule S3UploadRule \
    --event-bus-name "${EVENT_BUS_NAME}" \
    --targets '[
        {
            "Id": "1",
            "Arn": "'"${QUEUE_ARN}"'"
        }
    ]'

echo "Adding SQS Queue Policy for EventBridge..."
QUEUE_URL="http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/${SQS_NOTIFICATIONS_QUEUE}"

awslocal sqs set-queue-attributes \
    --queue-url "${QUEUE_URL}" \
    --attributes '{
        "Policy": "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"events.amazonaws.com\"},\"Action\":\"sqs:SendMessage\",\"Resource\":\"'"${QUEUE_ARN}"'\"}]}"
    }'

echo "EventBridge setup complete: S3 → EventBridge → SQS"