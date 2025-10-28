from pydantic import BaseModel

class QuestionPayload(BaseModel):
    user_id: str
    question: str
