import boto3
import json
import time
import logging
from urllib.parse import unquote
from config import AWS_REGION, AWS_ACCESS_KEY, AWS_SECRET_KEY, SQS_ENDPOINT, QUEUE_URL
from api.dto.document_dto import UploadPayload
from services.document_service import process_document

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

sqs = boto3.client(
    'sqs',
    endpoint_url=SQS_ENDPOINT,
    region_name=AWS_REGION,
    aws_access_key_id=AWS_ACCESS_KEY,
    aws_secret_access_key=AWS_SECRET_KEY
)

def parse_s3_key(key: str, size: int) -> dict:
    key = unquote(key)
    
    parts = key.split('/')
    if len(parts) < 3:
        raise ValueError(f"Invalid S3 key format: {key}")
    
    user_id = parts[1]  # user_34e0dGYAAfhiVU5SGoQsYYbZ3ug
    
    # Split document_id and file_name: "9eba6d20....-Slobodan Obradovic.pdf"
    filename_part = parts[2]
    
    if len(filename_part) > 37 and filename_part[36] == '-':
        document_id = filename_part[:36]
        file_name = filename_part[37:]
    else:
        raise ValueError(f"Cannot parse document_id from: {filename_part}")
    
    return {
        "user_id": user_id,
        "document_id": document_id,
        "file_name": file_name,
        "s3_key": key,
        "file_size": size
    }

def poll_sqs():
    logger.info(f"Starting SQS polling on: {QUEUE_URL}")
    
    while True:
        try:
            response = sqs.receive_message(
                QueueUrl=QUEUE_URL,
                MaxNumberOfMessages=10,
                WaitTimeSeconds=20,
                MessageAttributeNames=['All']
            )
            
            messages = response.get('Messages', [])
            
            if messages:
                logger.info(f"📨 Received {len(messages)} messages")
            
            for message in messages:
                try:
                    body = json.loads(message['Body'])
                    
                    
                    if 'detail' in body:
                        s3_info = body['detail']
                        bucket = s3_info['bucket']['name']
                        key = s3_info['object']['key']
                        size = s3_info['object']['size']
                    
                    elif 'Records' in body:
                        s3_info = body['Records'][0]['s3']
                        bucket = s3_info['bucket']['name']
                        key = s3_info['object']['key']
                        size = s3_info['object']['size']
                    else:
                        logger.error(f"Unknown format: {body}")
                        continue
                    
                    logger.info(f"Processing: {bucket}/{key}")
                    
                    parsed_data = parse_s3_key(key, size)
                    logger.info(f" Parsed: user={parsed_data['user_id']}, doc={parsed_data['document_id']}, file={parsed_data['file_name']}")
                    
                    
                    payload = UploadPayload(
                        bucket_name=bucket,
                        **parsed_data  # user_id, document_id, file_name, s3_key, file_size
                    )
                    
                    process_document(payload)
                    
                    
                    sqs.delete_message(
                        QueueUrl=QUEUE_URL,
                        ReceiptHandle=message['ReceiptHandle']
                    )
                    
                    logger.info(f"Done: {parsed_data['file_name']}")
                    
                except Exception as e:
                    logger.error(f"Error: {str(e)}", exc_info=True)
                    
        except Exception as e:
            logger.error(f"SQS polling error: {str(e)}")
            time.sleep(5)

if __name__ == "__main__":
    poll_sqs()