#!/bin/bash
set -e

awslocal s3api put-bucket-cors --bucket ${AWS_S3_BUCKET} --cors-configuration '{
  "CORSRules": [
    {
      "AllowedHeaders": ["*"],
      "AllowedMethods": ["GET", "PUT", "POST", "DELETE"],
      "AllowedOrigins": ["*"]
    }
  ]
}'
