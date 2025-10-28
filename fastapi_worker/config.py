import os

# AWS 
AWS_ENDPOINT = os.getenv("AWS_ENDPOINT", "http://localstack:4566")
AWS_REGION = os.getenv("AWS_REGION", "us-east-1")
AWS_ACCESS_KEY = os.getenv("AWS_ACCESS_KEY_ID", "test")
AWS_SECRET_KEY = os.getenv("AWS_SECRET_ACCESS_KEY", "test")

BUCKET_NAME = os.getenv("BUCKET_NAME", "my-documents-bucket")
DOWNLOAD_DIR = os.getenv("DOWNLOAD_DIR", "/tmp")


SQS_ENDPOINT = os.getenv("SQS_ENDPOINT", "http://localstack:4566")
QUEUE_URL = os.getenv(
    "QUEUE_URL", 
    "http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/notifications-queue"
)

# LLM
API_KEY = os.getenv("API_KEY", None) #Gemini
EMBEDDING_MODEL = os.getenv("EMBEDDING_MODEL", "models/text-embedding-004")
CHAT_MODEL = os.getenv("CHAT_MODEL", "gpt-4.1-mini")

# PINECONE 
PINECONE_API = os.getenv("PINECONE_API", None)
PINECONE_ENV = os.getenv("PINECONE_ENV", None)
PINECONE_INDEX = os.getenv("PINECONE_INDEX", "my-index")

# NEON DB
NEON_DB_URL = os.getenv("DB_URL", None)

