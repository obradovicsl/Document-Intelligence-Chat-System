from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import boto3
import os
from config import AWS_ENDPOINT, AWS_REGION, AWS_ACCESS_KEY, AWS_SECRET_KEY, BUCKET_NAME, DOWNLOAD_DIR

app = FastAPI()

# S3 client za LocalStack
s3 = boto3.client(
    "s3",
    endpoint_url=AWS_ENDPOINT,
    region_name=AWS_REGION,
    aws_access_key_id=AWS_ACCESS_KEY,
    aws_secret_access_key=AWS_SECRET_KEY,
)

# Payload model
class UploadPayload(BaseModel):
    user_id: str
    document_id: str
    file_name: str
    s3_key: str
    file_size: int

@app.post("/process-document")
async def process_document(payload: UploadPayload):
    try:
        local_path = os.path.join(DOWNLOAD_DIR, payload.document_id or "temp_file")
        s3.download_file(BUCKET_NAME, payload.s3_key, local_path)
        return {"status": "success", "message": f"File downloaded to {local_path}"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
