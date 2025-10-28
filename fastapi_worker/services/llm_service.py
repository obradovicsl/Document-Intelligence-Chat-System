import google.generativeai as genai
from config import EMBEDDING_MODEL, API_KEY, CHAT_MODEL

genai.configure(api_key=API_KEY)

def embed_texts(chunks: list[str]):
    embeddings = []
    for chunk in chunks:
        result = genai.embed_content(
            model=EMBEDDING_MODEL,
            content=chunk,
            task_type="retrieval_document"
        )
        embeddings.append(result['embedding'])
    return embeddings


def ask_llm(context_vector, payload):
    try:
        context = get_context(context_vector)

        prompt = f"""
        Odgovori na pitanje korisnika koristeći kontekst ispod.
        Ako odgovor ne postoji u kontekstu, reci da ne znaš.

        KONTEKST:
        {context}

        PITANJE:
        {payload.question}
        """

        print(prompt)

        try:
            model = genai.GenerativeModel(CHAT_MODEL)
            response = model.generate_content(prompt)
            
        except Exception as e:
            return "Došlo je do greške prilikom generisanja odgovora."
        try:
            answer = response.text.strip() if response.text else "Nema odgovora."
        except Exception as e:
            answer = "Greška pri obradi odgovora."

        return answer
    except Exception as e:
        return "Doslo je do greske prilikom obrade pitanja" + str(e)
    

def get_context(context_vector):
    context_parts = []
    
    for match in context_vector.get("matches", []):
        try:
            # Pinecone vraća dict, ne ScoredVector objekat
            if isinstance(match, dict):
                metadata = match.get("metadata", {})
            else:
                # fallback za objekat
                metadata = getattr(match, "metadata", {})
            
            # izvuci chunk_text
            text = metadata.get("chunk_text", "") if metadata else ""
            
            # osiguraj da je string
            if text and not isinstance(text, str):
                text = str(text)
            
            if text:  # dodaj samo ako ima teksta
                context_parts.append(text)
                
        except Exception as e:
            print(f"Greška pri parsiranju match-a: {e}")
            continue
    
    result = "\n\n".join(context_parts)
    print(f"Ekstrahovano {len(context_parts)} chunk-ova, ukupno {len(result)} karaktera")
    return result