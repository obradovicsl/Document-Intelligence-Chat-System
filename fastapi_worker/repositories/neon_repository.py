import psycopg2
from config import NEON_DB_URL
import uuid
from datetime import datetime

def save_chat_history(user_id, question, answer):
    chat_id = str(uuid.uuid4())  # generiši UUID kao string
    asked_at = datetime.now()    # trenutni datum/vreme
    with psycopg2.connect(NEON_DB_URL) as conn:
        with conn.cursor() as cur:
            cur.execute(
                "INSERT INTO chats (id, user_id, question, answer, asked_at) VALUES (%s, %s, %s, %s, %s)",
                (chat_id, user_id, question, answer, asked_at)
            )
