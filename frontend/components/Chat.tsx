// components/Chat.tsx
"use client";

import { useEffect, useRef } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Message } from "@/lib/types/chat.types";
import { Loader2, Trash2 } from "lucide-react";

interface ChatProps {
  messages: Message[];
  input: string;
  isLoading: boolean;
  onInputChange: (value: string) => void;
  onSubmit: () => void;
  onClear?: () => void;
}

export function Chat({
  messages,
  input,
  isLoading,
  onInputChange,
  onSubmit,
  onClear,
}: ChatProps) {
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      onSubmit();
    }
  };

  return (
    <Card className="w-full h-full flex flex-col shadow-lg">
      {/* Header */}
      <CardHeader className="border-b flex flex-row items-center justify-between">
        <CardTitle className="text-xl">Chat</CardTitle>
        {onClear && (
          <Button
            variant="ghost"
            size="sm"
            onClick={onClear}
            disabled={messages.length === 0}
          >
            <Trash2 className="h-4 w-4 mr-2" />
            Clear
          </Button>
        )}
      </CardHeader>

      {/* Messages */}
      <CardContent className="flex-1 p-0 overflow-hidden">
        <ScrollArea className="h-full p-4">
          <div ref={scrollRef} className="space-y-4">
            {messages.length === 0 && (
              <div className="text-center text-muted-foreground py-8">
                <p className="text-lg font-medium">No messages yet</p>
                <p className="text-sm">Ask a question about your documents</p>
              </div>
            )}

            {messages.map((message) => (
              <div
                key={message.id}
                className={`flex ${
                  message.role === "user" ? "justify-end" : "justify-start"
                }`}
              >
                <div
                  className={`max-w-[80%] rounded-lg p-4 ${
                    message.role === "user"
                      ? "bg-primary text-primary-foreground"
                      : "bg-muted"
                  }`}
                >
                  <div className="text-xs font-semibold mb-1 opacity-70">
                    {message.role === "user" ? "You" : "Assistant"}
                  </div>
                  <div className="whitespace-pre-wrap text-sm">
                    {message.content}
                  </div>
                  {message.timestamp && (
                    <div className="text-xs opacity-50 mt-2">
                      {new Date(message.timestamp).toLocaleTimeString()}
                    </div>
                  )}
                </div>
              </div>
            ))}

            {/* Loading Indicator */}
            {isLoading && (
              <div className="flex justify-start">
                <div className="bg-muted rounded-lg p-4">
                  <div className="flex items-center space-x-2">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    <span className="text-sm text-muted-foreground">
                      Thinking...
                    </span>
                  </div>
                </div>
              </div>
            )}
          </div>
        </ScrollArea>
      </CardContent>

      {/* Input */}
      <div className="p-4 border-t">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            onSubmit();
          }}
          className="flex gap-2"
        >
          <Input
            value={input}
            onChange={(e) => onInputChange(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Ask a question..."
            disabled={isLoading}
            className="flex-1"
          />
          <Button type="submit" disabled={isLoading || !input.trim()}>
            {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : "Send"}
          </Button>
        </form>
        <p className="text-xs text-muted-foreground mt-2">
          Press Enter to send, Shift+Enter for new line
        </p>
      </div>
    </Card>
  );
}
