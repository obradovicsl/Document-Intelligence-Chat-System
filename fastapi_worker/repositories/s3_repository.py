import boto3
from config import AWS_REGION, AWS_ACCESS_KEY, AWS_SECRET_KEY, BUCKET_NAME, AWS_ENDPOINT

s3 = boto3.client(
    "s3",
    endpoint_url=AWS_ENDPOINT,
    region_name=AWS_REGION,
    aws_access_key_id=AWS_ACCESS_KEY,
    aws_secret_access_key=AWS_SECRET_KEY,
)

def download_file(s3_key: str) -> bytes:
    obj = s3.get_object(Bucket=BUCKET_NAME, Key=s3_key)
    return obj['Body'].read()
