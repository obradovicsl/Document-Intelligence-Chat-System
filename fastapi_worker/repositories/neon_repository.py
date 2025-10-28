import psycopg2
from config import NEON_DB_URL

def save_chat_history(user_id, question, answer):
    with psycopg2.connect(NEON_DB_URL) as conn:
        with conn.cursor() as cur:
            cur.execute(
                "INSERT INTO chat_history (user_id, question, answer) VALUES (%s, %s, %s)",
                (user_id, question, answer)
            )
