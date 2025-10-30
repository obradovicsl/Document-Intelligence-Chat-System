# RAG-Powered Document Intelligence Chat System

A modern document intelligence platform that enables users to upload documents and interact with them through natural language conversations, powered by Retrieval-Augmented Generation (RAG).

## Overview

Upload your PDFs and text documents, and ask questions in natural language. The system automatically processes documents, creates semantic embeddings, and uses AI to provide accurate answers based on your document content.

## Architecture
![RAG Architecture Diagram](diagrams/RAG%20Dijagram.png)

## Tech Stack

**Frontend**
- Next.js 14 (App Router)
- shadcn/ui components
- Clerk authentication
- TypeScript

**Backend**
- Go 1.21+ (GORM, AWS SDK v2, pgx)
- Python 3.12 (FastAPI, boto3, PyPDF2, LangChain)

**AI/ML**
- Google Gemini
  - text-embedding-004 (embeddings)
  - gemini-2.5-flash-lite (chat completions)
- Pinecone (vector database)

**Infrastructure**
- Docker & Docker Compose
- LocalStack (S3, SQS, EventBridge)
- NeonDB (PostgreSQL)

## Prerequisites

- Docker & Docker Compose installed
- Node.js 18+ and npm
- Go 1.21+
- Python 3.12+
- Gemini API key
- Pinecone account and API key
- Clerk account for authentication
- NeonDB PostgreSQL database

## Document Processing Flow

1. **Upload Initiation**
   - User selects document in frontend
   - Frontend requests presigned S3 URL from Go API
   - Go API generates URL and creates document record in NeonDB

2. **File Upload**
   - Frontend uploads file directly to S3 using presigned URL
   - S3 key format: `documents/{user_id}/{document_id}-{filename}`

3. **Event Triggering**
   - S3 PutObject event triggers EventBridge rule
   - EventBridge routes event to SQS queue

4. **Async Processing**
   - Python worker polls SQS queue (long polling, 20s)
   - Worker parses S3 event and extracts metadata
   - Downloads document from S3
   - Extracts text from PDF/document
   - Chunks text using LangChain text splitter
   - Generates embeddings for each chunk
   - Stores vectors in Pinecone with metadata
   - Updates document status in NeonDB
   - Deletes message from SQS queue

5. **Ready for Chat**
   - Document status changes to "completed"
   - User can now ask questions about the document

## Chat Query Flow

1. **Question Submission**
   - User types question in chat interface
   - Frontend sends question to Go API with document ID

2. **Query Processing**
   - Go API forwards request to Python RAG service
   - RAG service embeds the question using text-embedding-3-small

3. **Semantic Search**
   - Embedded query searches Pinecone for similar vectors
   - Top-k most relevant document chunks retrieved
   - Chunks filtered by document ID

4. **Context Retrieval**
   - Document metadata fetched from NeonDB
   - Relevant chunks combined into context

5. **Answer Generation**
   - Context + question sent to GPT-4.1-mini
   - Model generates answer based on document content
   - Answer includes source references

6. **Response & Storage**
   - Answer returned to frontend
   - Question-answer pair saved to chat_history table
   - Chat history displayed in UI
