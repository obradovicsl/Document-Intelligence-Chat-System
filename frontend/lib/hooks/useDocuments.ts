// hooks/useDocuments.ts
import { useQuery } from "@tanstack/react-query";
import { getUserDocuments } from "@/lib/services/document.service";
import { DocumentListResponse } from "@/lib/types/document.types";
import { useEffect, useState } from "react";
import { useAuth } from "@clerk/nextjs";

export const useDocuments = () => {
  const [data, setData] = useState<DocumentListResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const { getToken } = useAuth();

  useEffect(() => {
    const fetchDocuments = async () => {
      try {
        const token = await getToken();
        if (!token) throw new Error("User not authenticated");

        const res = await getUserDocuments(token);
        setData(res);
      } catch (err) {
        setError(err as Error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchDocuments();
  }, []);

  return {
    data,
    isLoading,
    isError: !!error,
    error,
    documents: data?.documents ?? [],
    count: data?.count ?? 0,
  };
};
