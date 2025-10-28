// types/document.ts
export interface DocumentDTO {
  iD: string;
  userID: string;
  fileName: string;
  fileSize: number;
  status: string;
  createdAt: string;
}

export interface DocumentListResponse {
  documents: DocumentDTO[];
  count: number;
}
