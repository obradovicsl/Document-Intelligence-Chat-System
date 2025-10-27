import os
from io import BytesIO
import psycopg2
import boto3
from PyPDF2 import PdfReader
from openai import OpenAI
from pinecone import Pinecone, ServerlessSpec
from config import AWS_ENDPOINT, AWS_REGION, AWS_ACCESS_KEY, AWS_SECRET_KEY, BUCKET_NAME
from config import PINECONE_API, PINECONE_ENV, EMBEDDING_MODEL, API_KEY


# ----- S3 -----
s3 = boto3.client(
    "s3",
    endpoint_url=AWS_ENDPOINT,
    region_name=AWS_REGION,
    aws_access_key_id=AWS_ACCESS_KEY,
    aws_secret_access_key=AWS_SECRET_KEY,
)

def download_from_s3(s3_key: str) -> bytes:
    try:
        print(s3_key)
        obj = s3.get_object(Bucket=BUCKET_NAME, Key=s3_key)
        return obj['Body'].read()
    except Exception as e:
        print("Error downloading from S3:", str(e))
        raise


# ----- PDF Parsing -----
def parse_pdf(file_bytes: bytes) -> str:
    reader = PdfReader(BytesIO(file_bytes))
    text = "".join(page.extract_text() for page in reader.pages)
    return text


# ----- Chunking -----
def chunk_text(text: str, chunk_size=1000, overlap=200):
    chunks = []
    start = 0
    text_len = len(text)
    
    while start < text_len:
        end = start + chunk_size
        chunk = text[start:end]
        chunks.append(chunk)
        start += chunk_size - overlap
    return chunks



# ----- Pinecone setup -----
pc = Pinecone(api_key=PINECONE_API)
index_name = "my-index"

if index_name not in pc.list_indexes().names():
    pc.create_index(
        name=index_name,
        dimension=1536,  # embedding dimension
        metric="cosine",
        spec=ServerlessSpec(
            cloud="aws",
            region="us-east-1"
        ),
    )

index = pc.Index(index_name)


def embed_and_upsert(chunks: list[str], metadata: dict):
    client = OpenAI(api_key=API_KEY)
    vectors_to_upsert = []
    
    try:
        for i, chunk in enumerate(chunks):
            print(f"Embedding chunk {i}/{len(chunks)}")
            resp = client.embeddings.create(model=EMBEDDING_MODEL, input=chunk)
            vectors_to_upsert.append({
                "id": f"{metadata['document_id']}_chunk_{i}",
                "values": resp.data[0].embedding,
                "metadata": {**metadata, "chunk_index": i}
            })
        
        print("Upserting vectors to Pinecone...")
        index.upsert(vectors=vectors_to_upsert)
        print("Upsert successful")

    except Exception as e:
        print("Error in embed_and_upsert:", e)
        raise e