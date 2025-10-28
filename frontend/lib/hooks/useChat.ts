// lib/hooks/useChat.ts
"use client";

import { useState, useCallback } from "react";
import { chatService } from "@/lib/services/chat.service";
import { Message, ChatError } from "@/lib/types/chat.types";
import { useAuth } from "@clerk/nextjs";

interface UseChatOptions {
  userId: string;
  onError?: (error: ChatError) => void;
}

export function useChat({ userId, onError }: UseChatOptions) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const { getToken } = useAuth();

  const sendMessage = useCallback(
    async (content?: string) => {
      const messageContent = content || input.trim();
      if (!messageContent) return;

      const userMessage: Message = {
        id: `user-${Date.now()}`,
        role: "user",
        content: messageContent,
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, userMessage]);
      setInput("");
      setIsLoading(true);

      try {
        const token = await getToken();
        if (!token) {
          throw new Error("User not authenticated");
        }

        const response = await chatService.sendMessage(
          {
            question: messageContent,
            userId,
          },
          token
        );

        const assistantMessage: Message = {
          id: `assistant-${Date.now()}`,
          role: "assistant",
          content: response.answer,
          timestamp: new Date(),
        };

        setMessages((prev) => [...prev, assistantMessage]);
      } catch (error) {
        const chatError: ChatError = {
          message: error instanceof Error ? error.message : "Unknown error",
        };

        const errorMessage: Message = {
          id: `error-${Date.now()}`,
          role: "assistant",
          content: `Sorry, there was an error: ${chatError.message}`,
          timestamp: new Date(),
        };

        setMessages((prev) => [...prev, errorMessage]);

        if (onError) {
          onError(chatError);
        }
      } finally {
        setIsLoading(false);
      }
    },
    [input, userId, onError]
  );

  const clearChat = useCallback(() => {
    setMessages([]);
    setInput("");
  }, []);

  const removeMessage = useCallback((messageId: string) => {
    setMessages((prev) => prev.filter((m) => m.id !== messageId));
  }, []);

  return {
    messages,
    input,
    setInput,
    isLoading,
    sendMessage,
    clearChat,
    removeMessage,
  };
}
