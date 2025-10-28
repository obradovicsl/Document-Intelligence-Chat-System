// components/DocumentList.tsx
"use client";

import { Document } from "./Document";
import { useDocuments } from "@/lib/hooks/useDocuments";

export const DocumentList = () => {
  const { data, isLoading, isError, error, documents, count } = useDocuments();

  if (isLoading) return <p>Loading...</p>;
  if (isError) return <div>Error: {error?.message}</div>;
  if (!documents || documents.length === 0)
    return <div>No documents found.</div>;

  return (
    <div className="space-y-4">
      <h2 className="text-lg font-semibold">Documents ({count})</h2>
      {documents.map((doc) => {
        console.log(doc);
        return <Document key={doc.iD} document={doc} />;
      })}
    </div>
  );
};
