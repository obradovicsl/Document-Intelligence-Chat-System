// app/chat/page.tsx
"use client";

import { Chat } from "@/components/Chat";
import { useChat } from "@/lib/hooks/useChat";
import { useState } from "react";

export default function ChatPage() {
  const [userId] = useState("user_123");
  const [selectedDocumentId, setSelectedDocumentId] = useState<string>();

  const { messages, input, setInput, isLoading, sendMessage, clearChat } =
    useChat({
      userId,
      documentId: selectedDocumentId,
      onError: (error) => {
        console.error("Chat error:", error);
      },
    });

  return (
    <div className="flex h-screen">
      <main className="flex-1">
        <Chat
          messages={messages}
          input={input}
          isLoading={isLoading}
          onInputChange={setInput}
          onSubmit={() => sendMessage()}
          onClear={clearChat}
        />
      </main>
    </div>
  );
}
