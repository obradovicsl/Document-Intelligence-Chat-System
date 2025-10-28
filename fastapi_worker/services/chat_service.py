from services.llm_service import embed_texts, handle_question
from repositories.pinecone_repository import query_vectors
from repositories.neon_repository import save_chat_history

def handle_question(payload):
    embedding = embed_texts([payload.question])[0]

    # Pinecone query
    query_results = query_vectors(embedding, top_k=5)
    
    # LLM service
    answer = handle_question(query_results, payload)

    # Save to NeonDB
    save_chat_history(payload.user_id, payload.question, answer)

    return answer