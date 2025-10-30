import { useState } from "react";
import { useAuth } from "@clerk/nextjs";
import {
  requestPresignedURL,
  uploadToS3,
  UploadRequest,
  notifyAPI,
} from "@/lib/services/upload.service";

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

      const { uploadUrl, documentId, userId, key } = await requestPresignedURL(
        uploadData,
        token
      );

      console.log(uploadUrl);

      await uploadToS3(file, uploadUrl, {
        user_id: userId,
        document_id: documentId,
        file_name: file.name,
        s3_key: key,
        file_size: file.size,
      });

      // await notifyAPI(
      //   {
      //     user_id: userId,
      //     document_id: documentId,
      //     file_name: file.name,
      //     s3_key: key,
      //     file_size: file.size,
      //   },
      //   token
      // );
    } catch (err: any) {
      setError(err.message || "Upload failed");
    } finally {
      setLoading(false);
    }
  };

  return { uploadFile, loading, error };
}
