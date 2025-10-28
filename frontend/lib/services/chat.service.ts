import { ChatRequest, ChatResponse } from "@/lib/types/chat.types";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_GO_API_URL || "http://localhost:8000";

class ChatService {
  /**
   * Šalje pitanje na Go API i vraća odgovor
   */
  async sendMessage(request: ChatRequest): Promise<ChatResponse> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/question`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          // Dodaj auth token ako imaš
          // "Authorization": `Bearer ${getToken()}`
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(
          errorData.message || `HTTP error! status: ${response.status}`
        );
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error("Chat service error:", error);
      throw error;
    }
  }

  /**
   * Streaming verzija (ako Go API podržava)
   */
  async sendMessageStream(
    request: ChatRequest,
    onChunk: (chunk: string) => void
  ): Promise<void> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/question/stream`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body?.getReader();
      const decoder = new TextDecoder();

      if (!reader) {
        throw new Error("No response body");
      }

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value);
        onChunk(chunk);
      }
    } catch (error) {
      console.error("Streaming error:", error);
      throw error;
    }
  }
}

export const chatService = new ChatService();
