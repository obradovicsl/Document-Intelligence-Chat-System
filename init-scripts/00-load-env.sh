#!/bin/bash
set -e

echo "========================================="
echo "Loading and validating environment..."
echo "========================================="

validate_var() {
    local var_name=$1
    local var_value=${!var_name}
    
    if [ -z "$var_value" ]; then
        echo "❌ ERROR: $var_name is not set!"
        return 1
    else
        echo "✓ $var_name is set"
        return 0
    fi
}

echo ""
echo "Validating AWS configuration..."
validate_var "AWS_REGION"
validate_var "AWS_ENDPOINT"

echo ""
echo "Validating S3 buckets..."
validate_var "AWS_S3_BUCKET"

echo ""
echo "Validating SQS queues..."
validate_var "SQS_NOTIFICATIONS_QUEUE"

echo ""
echo "========================================="
echo "✓ All environment variables validated!"
echo "========================================="
echo ""

# Logging all env vars
echo "Current configuration:"
echo "  AWS Region: $AWS_REGION"
echo "  AWS Endpoint: $AWS_ENDPOINT"
echo "  S3 Bucket: $AWS_S3_BUCKET"
echo "  SQS Notification Queue: $SQS_NOTIFICATIONS_QUEUE"
echo ""