from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from utils import download_from_s3, parse_pdf, chunk_text, embed_and_upsert

app = FastAPI()


class UploadPayload(BaseModel):
    user_id: str
    document_id: str
    file_name: str
    s3_key: str
    file_size: int

@app.post("/process-document")
async def process_document(payload: UploadPayload):
    try:
        print(payload.document_id, payload.s3_key, payload.file_size, payload.file_name)
        print("Downloading file from S3")
        file_bytes = download_from_s3(payload.s3_key)
        
        print("Parsing PDF...")
        text = parse_pdf(file_bytes)
        
        print("Chunking text")
        chunks = chunk_text(text)

        metadata = {
            "document_id": payload.document_id,
            "user_id": payload.user_id,
            "file_name": payload.file_name,
            "s3_key": payload.s3_key,
            "file_size": payload.file_size
        }

        # Embedding
        print("Embedding and upserting")
        embed_and_upsert(chunks, metadata)

        print("DONE!")
        return {"status": "success", "message": "File chunked and embedded"}

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
