import { useState } from "react";
import { useAuth } from "@clerk/nextjs";
import {
  requestPresignedURL,
  uploadToS3,
  UploadRequest,
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

      const { uploadUrl } = await requestPresignedURL(uploadData, token);
      //   await uploadToS3(file, uploadUrl);
      console.log(uploadUrl);
    } catch (err: any) {
      setError(err.message || "Upload failed");
    } finally {
      setLoading(false);
    }
  };

  return { uploadFile, loading, error };
}
