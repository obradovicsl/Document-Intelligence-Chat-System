export type Message = {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp?: Date;
};

export type ChatRequest = {
  question: string;
  userId: string;
};

export type ChatResponse = {
  answer: string;
  sources?: Array<{
    documentId: string;
    fileName: string;
    chunkText: string;
    score: number;
  }>;
};

export type ChatError = {
  message: string;
  code?: string;
};
