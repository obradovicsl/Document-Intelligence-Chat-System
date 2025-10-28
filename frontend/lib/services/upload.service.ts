export const API_URL = process.env.NEXT_PUBLIC_API_URL;

export type UploadRequest = {
  fileName: string;
  fileType: string;
  fileSize: number;
};

export type UploadResponse = {
  uploadUrl: string;
  documentId: string;
  userId: string;
  key: string;
};

export interface UploadCompletePayload {
  user_id: string;
  document_id: string;
  file_name: string;
  s3_key: string;
  file_size: number;
}

export async function requestPresignedURL(
  uploadData: UploadRequest,
  token: string
): Promise<UploadResponse> {
  const res = await fetch(`${API_URL}/api/upload/init`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(uploadData),
  });

  if (!res.ok) throw new Error("Failed to request presigned URL");
  return res.json();
}

export async function uploadToS3(file: File, presignedUrl: string) {
  const res = await fetch(presignedUrl, {
    method: "PUT",
    body: file,
    headers: {
      "Content-Type": file.type,
    },
  });

  if (!res.ok) throw new Error("Failed to upload file to S3");
}

export async function notifyAPI(payload: UploadCompletePayload, token: string) {
  const res = await fetch(`${API_URL}/api/upload/complete`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Failed to notify API: ${res.status} ${text}`);
  }

  return await res.json(); // ili možeš vratiti status OK ako backend ne vraća JSON
}
