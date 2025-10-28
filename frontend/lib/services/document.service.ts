import { DocumentListResponse } from "@/lib/types/document.types";

export const API_URL = process.env.NEXT_PUBLIC_API_URL;

export const getUserDocuments = async (
  token: string
): Promise<DocumentListResponse> => {
  //   const res = await fetch(`${API_URL}/api/documents/user/me`);

  const res = await fetch(`${API_URL}/api/documents/user/me`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) throw new Error("Failed to fetch documents");

  const data: DocumentListResponse = await res.json();
  return data;
};
