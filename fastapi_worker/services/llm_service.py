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


def handle_question(context_vector, payload):
    context = "\n\n".join([match['metadata']['chunk_text'] for match in context_vector['matches']])

    prompt = f"""
    Odgovori na pitanje korisnika koristeći kontekst ispod.
    Ako odgovor ne postoji u kontekstu, reci da ne znaš.

    KONTEKST:
    {context}

    PITANJE:
    {payload.question}
    """

    model = genai.GenerativeModel(CHAT_MODEL)
    response = model.generate_content(prompt)

    answer = response.text.strip() if response.text else "Nema odgovora."

    return answer