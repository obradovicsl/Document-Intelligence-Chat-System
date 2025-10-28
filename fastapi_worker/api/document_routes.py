from fastapi import APIRouter, HTTPException
from api.dto.document_dto import UploadPayload
from services.document_service import process_document

router = APIRouter()

@router.post("/process")
async def process_uploaded_document(payload: UploadPayload):
    try:
        process_document(payload)
        return {"status": "success", "message": "Document processed and embedded successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
