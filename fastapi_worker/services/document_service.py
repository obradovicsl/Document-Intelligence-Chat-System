from repositories.s3_repository import download_file
from repositories.pinecone_repository import upsert_vectors
from utils import parse_pdf
from utils import chunk_text
from services.llm_service import embed_texts

def process_document(payload):
    file_bytes = download_file(payload.s3_key)
    text = parse_pdf(file_bytes)
    chunks = chunk_text(text)

    embeddings = embed_texts(chunks)
    vectors = [
        {
            "id": f"{payload.document_id}_chunk_{i}",
            "values": embeddings[i],
            "metadata": {
                **payload.dict(),
                "chunk_index": i,
                "chunk_text": chunk,
            },
        }
        for i, chunk in enumerate(chunks)
    ]

    upsert_vectors(vectors)
