// components/Document.tsx
"use client";

import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { DocumentDTO } from "@/lib/types/document.types";
import { format } from "date-fns";

interface DocumentProps {
  document: DocumentDTO;
}

export const Document = ({ document }: DocumentProps) => {
  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle>{document.fileName}</CardTitle>
        <CardDescription>
          Status: {document.status} • Size: {document.fileSize} bytes
        </CardDescription>
      </CardHeader>
      <CardContent>
        Created at: {format(new Date(document.createdAt), "PPpp")}
      </CardContent>
    </Card>
  );
};
