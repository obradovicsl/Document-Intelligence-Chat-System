import os
from io import BytesIO
import psycopg2
import boto3
from PyPDF2 import PdfReader
import google.generativeai as genai
from pinecone import Pinecone, ServerlessSpec
from config import (
    AWS_ENDPOINT, 
    AWS_REGION, 
    AWS_ACCESS_KEY, 
    AWS_SECRET_KEY, 
    BUCKET_NAME,
    PINECONE_API, 
    PINECONE_ENV, 
    EMBEDDING_MODEL, 
    API_KEY
)


genai.configure(api_key=API_KEY)

# ----- S3 -----
s3 = boto3.client(
    "s3",
    endpoint_url=AWS_ENDPOINT,
    region_name=AWS_REGION,
    aws_access_key_id=AWS_ACCESS_KEY,
    aws_secret_access_key=AWS_SECRET_KEY,
)

def download_from_s3(s3_key: str) -> bytes:
    """Download file from S3 bucket"""
    try:
        print(f"Downloading from S3: {s3_key}")
        obj = s3.get_object(Bucket=BUCKET_NAME, Key=s3_key)
        return obj['Body'].read()
    except Exception as e:
        print(f"Error downloading from S3: {str(e)}")
        raise

# ----- PDF Parsing -----
def parse_pdf(file_bytes: bytes) -> str:
    """Extract text from PDF bytes"""
    try:
        reader = PdfReader(BytesIO(file_bytes))
        text = ""
        for page in reader.pages:
            text += page.extract_text()
        return text
    except Exception as e:
        print(f"Error parsing PDF: {str(e)}")
        raise

# ----- Chunking -----
def chunk_text(text: str, chunk_size: int = 1000, overlap: int = 200) -> list[str]:
    """Split text into overlapping chunks"""
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
        dimension=768,  # embedding dimension
        metric="cosine",
        spec=ServerlessSpec(
            cloud="aws",
            region="us-east-1"
        ),
    )

index = pc.Index(index_name)

# ----- Embedding & Upsert -----
def embed_and_upsert(chunks: list[str], metadata: dict):
    """Generate embeddings and upsert to Pinecone"""
    print("Starting embedding and upsert process...")
    vectors_to_upsert = []
    
    try:
        # Process each chunk
        for i, chunk in enumerate(chunks):
            print(f"Embedding chunk {i+1}/{len(chunks)}")
            
            # Generate embedding using Gemini
            result = genai.embed_content(
                model=EMBEDDING_MODEL,
                content=chunk,
                task_type="retrieval_document"
            )
            
            embedding = result['embedding']
            
            # Prepare vector for upsert
            vector_id = f"{metadata['document_id']}_chunk_{i}"
            vectors_to_upsert.append({
                "id": vector_id,
                "values": embedding,
                "metadata": {
                    **metadata, 
                    "chunk_index": i,
                    "chunk_text": chunk[:500]  # Store first 500 chars for reference
                }
            })
        
        # Upsert all vectors to Pinecone
        print(f"Upserting {len(vectors_to_upsert)} vectors to Pinecone...")
        index.upsert(vectors=vectors_to_upsert)
        print("Upsert successful!")
        
        return True
        
    except Exception as e:
        print(f"Error in embed_and_upsert: {str(e)}")
        raise e