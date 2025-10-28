from fastapi import FastAPI
from api import document_routes, chat_routes

app = FastAPI(title="Python Processing Service")

app.include_router(document_routes.router, prefix="/documents", tags=["Documents"])
app.include_router(chat_routes.router, prefix="/chat", tags=["Chat"])