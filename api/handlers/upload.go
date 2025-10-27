package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/auth"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/services"
)

type UploadRequest struct {
    FileName    string `json:"fileName"`
    FileType    string `json:"fileType"`
    FileSize    int64  `json:"fileSize"`
}

type UploadResponse struct {
    UploadURL  string `json:"uploadUrl"`
    DocumentID string `json:"documentId"`
    Key        string `json:"key"`
}

func GeneratePresignedURL(w http.ResponseWriter, r *http.Request) {
    userID := auth.GetUserID(r)

    var req UploadRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate file
    if req.FileName == "" {
        http.Error(w, "File name is required", http.StatusBadRequest)
        return
    }

    // Generate unique key for S3
    key := services.GenerateS3Key(userID, req.FileName)

    // Generate presigned URL
    uploadURL, err := services.GeneratePresignedUploadURL(key, req.FileType, 15*time.Minute)
    if err != nil {
        http.Error(w, "Failed to generate upload URL: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Create document record in database
    documentID, err := services.CreateDocumentRecord(userID, req.FileName, key, req.FileSize)
    if err != nil {
        http.Error(w, "Failed to create document record: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Return response
    response := UploadResponse{
        UploadURL:  uploadURL,
        DocumentID: documentID,
        Key:        key,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}