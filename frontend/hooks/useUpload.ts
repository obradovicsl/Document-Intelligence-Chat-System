import { useState } from "react";
import { useAuth } from "@clerk/nextjs";
import {
  requestPresignedURL,
  uploadToS3,
  UploadRequest,
  notifyAPI,
} from "@/services/uploadService";

export function useUpload() {
  const { getToken } = useAuth();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const uploadFile = async (file: File) => {
    setLoading(true);
    setError(null);

    try {
      // Clerk JWT
      const token = await getToken();

      if (!token) {
        throw new Error("User not authenticated");
      }

      console.log(token);

      const uploadData: UploadRequest = {
        fileName: file.name,
        fileType: file.type,
        fileSize: file.size,
      };

      const { uploadUrl, documentId, userId } = await requestPresignedURL(
        uploadData,
        token
      );
      await uploadToS3(file, uploadUrl);

      await notifyAPI({
        user_id: userId,
        document_id: documentId,
        file_name: file.name,
        s3_key: uploadUrl.split(".com/")[1],
        file_size: file.size,
      });
    } catch (err: any) {
      setError(err.message || "Upload failed");
    } finally {
      setLoading(false);
    }
  };

  return { uploadFile, loading, error };
}
