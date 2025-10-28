"use client";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useUpload } from "@/lib/hooks/useUpload";

export default function UploadBox() {
  const [file, setFile] = useState<File | null>(null);
  const { uploadFile, loading, error } = useUpload();

  const handleUpload = () => {
    if (!file) return;
    uploadFile(file);
  };

  return (
    <div className="p-4 border rounded-xl">
      <Input
        type="file"
        onChange={(e) => setFile(e.target.files?.[0] || null)}
      />
      <Button onClick={handleUpload} className="mt-2" disabled={loading}>
        {loading ? "Uploading..." : "Upload"}
      </Button>
      {error && <p className="text-red-500 mt-2">{error}</p>}
    </div>
  );
}
