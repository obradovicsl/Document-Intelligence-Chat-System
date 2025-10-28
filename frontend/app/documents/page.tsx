"use client";
import { DocumentList } from "@/components/DocumentList";

export default function DocumentsPage() {
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">My Documents</h1>
      <DocumentList />
    </div>
  );
}
