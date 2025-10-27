import os

AWS_ENDPOINT = os.getenv("AWS_ENDPOINT", "http://localhost:4566")
AWS_REGION = os.getenv("AWS_REGION", "us-east-1")
AWS_ACCESS_KEY = os.getenv("AWS_ACCESS_KEY_ID", "test")
AWS_SECRET_KEY = os.getenv("AWS_SECRET_ACCESS_KEY", "test")
BUCKET_NAME = os.getenv("BUCKET_NAME", "my-documents-bucket")
DOWNLOAD_DIR = os.getenv("DOWNLOAD_DIR", "/tmp")
API_KEY = os.getenv("API_KEY", None)
EMBEDDING_MODEL = os.getenv("EMBEDDING_MODEL", "text-embedding-3-small")
CHAT_MODEL = os.getenv("CHAT_MODEL", "gpt-4.1-mini")
PINECONE_API = os.getenv("PINECONE_API", None)
PINECONE_ENV = os.getenv("PINECONE_ENV", None)