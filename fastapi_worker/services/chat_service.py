from services.llm_service import embed_texts, ask_llm
from repositories.pinecone_repository import query_vectors
from repositories.neon_repository import save_chat_history

def handle_question(payload):
    try:
        print("embedding")
        embedding = embed_texts([payload.question])[0]
        print(embedding)
        print("querying")
        # Pinecone query
        query_results = query_vectors(embedding, top_k=5)
        print(query_results)
        
        print("handling question")
        # LLM service
        answer = ask_llm(query_results, payload)

        print("saving to db")
        # Save to NeonDB
        save_chat_history(payload.user_id, payload.question, answer)
    
        return answer
    
    except Exception as e:
        return "Doslo je do greske " + str(e)