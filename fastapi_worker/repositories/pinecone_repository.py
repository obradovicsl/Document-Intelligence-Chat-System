from pinecone import Pinecone, ServerlessSpec
from config import PINECONE_API, PINECONE_INDEX

pc = Pinecone(api_key=PINECONE_API)

index_name = PINECONE_INDEX

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

def upsert_vectors(vectors):
    index.upsert(vectors=vectors)

def query_vectors(embedding, top_k=10, filter=None):
    return index.query(vector=embedding, top_k=top_k, filter=filter, include_metadata=True)