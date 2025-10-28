from fastapi import APIRouter, HTTPException
from api.dto.chat_dto import QuestionPayload
from services.chat_service import handle_question

router = APIRouter()

@router.post("/ask")
async def ask_question(payload: QuestionPayload):
    try:
        answer = handle_question(payload)
        return {"status": "success", "answer": answer}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
