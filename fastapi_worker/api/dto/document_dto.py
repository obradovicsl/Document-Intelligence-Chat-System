from pydantic import BaseModel

class UploadPayload(BaseModel):
    user_id: str
    document_id: str
    file_name: str
    s3_key: str
    file_size: int
