#!/bin/bash
set -euo pipefail

echo "Creating simple SQS queues..."

create_queue_if_not_exists() {
    local queue_name=$1

    if ! awslocal sqs get-queue-url --queue-name "${queue_name}" >/dev/null 2>&1; then
        echo "Creating queue: ${queue_name}"
        awslocal sqs create-queue --queue-name "${queue_name}"
        echo "✓ Queue ${queue_name} created"
    else
        echo "⊘ Queue ${queue_name} already exists"
    fi
}

create_queue_if_not_exists "${SQS_NOTIFICATIONS_QUEUE}"

echo "SQS queues setup complete!"
