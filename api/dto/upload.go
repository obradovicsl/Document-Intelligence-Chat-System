package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/models"
)

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

type DocumentDTO struct {
    ID        uuid.UUID `json:"id"`
    UserID    string    `json:"userId"`
    FileName  string    `json:"fileName"`
    FileSize  int64     `json:"fileSize"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
}

type DocumentListResponse struct {
    Documents []DocumentDTO `json:"documents"`
    Count     int           `json:"count"`
}


func ToDocumentDTO(d *models.Document) DocumentDTO {
    return DocumentDTO{
        ID:        d.ID,
        UserID:    d.UserID,
        FileName:  d.FileName,
        FileSize:  d.FileSize,
        Status:    string(d.Status),
        CreatedAt: d.CreatedAt,
    }
}

func ToDocumentDTOList(docs []*models.Document) []DocumentDTO {
    dtoList := make([]DocumentDTO, len(docs))
    for i, d := range docs {
        dtoList[i] = ToDocumentDTO(d)
    }
    return dtoList
}