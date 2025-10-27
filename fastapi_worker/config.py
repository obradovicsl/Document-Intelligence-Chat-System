import os

AWS_ENDPOINT = os.getenv("AWS_ENDPOINT", "http://localhost:4566")
AWS_REGION = os.getenv("AWS_REGION", "us-east-1")
AWS_ACCESS_KEY = os.getenv("AWS_ACCESS_KEY_ID", "test")
AWS_SECRET_KEY = os.getenv("AWS_SECRET_ACCESS_KEY", "test")
BUCKET_NAME = os.getenv("BUCKET_NAME", "my-documents-bucket")
DOWNLOAD_DIR = os.getenv("DOWNLOAD_DIR", "/tmp")