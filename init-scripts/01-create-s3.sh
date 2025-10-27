#!/bin/bash
set -e

echo "Creating S3 buckets..."

create_bucket_if_not_exists() {
    local bucket_name=$1
    
    if awslocal s3 ls "s3://${bucket_name}" 2>&1 | grep -q 'NoSuchBucket'; then
        echo "Creating bucket: ${bucket_name}"
        awslocal s3 mb "s3://${bucket_name}"
        echo "✓ Bucket ${bucket_name} created"
    else
        echo "⊘ Bucket ${bucket_name} already exists"
    fi
}

create_bucket_if_not_exists "${AWS_S3_BUCKET}"

echo "S3 buckets setup complete!"