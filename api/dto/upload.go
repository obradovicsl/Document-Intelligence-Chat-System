package dto

type UploadDocumentRequest struct {
    FileName    string `json:"fileName"`
    FileType    string `json:"fileType"`
    FileSize    int64  `json:"fileSize"`
}

type UploadDocumentResponse struct {
    UploadURL  string `json:"uploadUrl"`
    DocumentID string `json:"documentId"`
    UserID string `json:"userId"`
    Key        string `json:"key"`
}

type UploadDocumentPayload struct {
    UserID     string `json:"user_id"`
    DocumentID string `json:"document_id"`
    FileName   string `json:"file_name"`
    S3Key      string `json:"s3_key"`
    FileSize   int64  `json:"file_size"`
}

