package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/auth"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/models"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/services"
)

func GeneratePresignedURL(w http.ResponseWriter, r *http.Request) {
    userID := auth.GetUserID(r)

    slog.Info("decoding request")
    var req models.UploadRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        slog.Error("Invalid request body")
        return
    }

    // Validate file
    if req.FileName == "" {
        http.Error(w, "File name is required", http.StatusBadRequest)
        slog.Error("file name is required")
        return
    }

    // Generate unique key for S3
    slog.Info("generating s3 key")
    key := services.GenerateS3Key(userID, req.FileName)
    
    // Generate presigned URL
    slog.Info("generating presigned URL")
    uploadURL, err := services.GeneratePresignedUploadURL(key, req.FileType, 5*time.Minute)
    if err != nil {
        http.Error(w, "Failed to generate upload URL: "+err.Error(), http.StatusInternalServerError)
        slog.Error("failed to generate upload URL", "error", err.Error())
        return
    }
    
    // Create document record in database
    slog.Info("creating document record in database")
    documentID, err := services.CreateDocumentRecord(userID, req.FileName, key, req.FileSize)
    if err != nil {
        slog.Error("failed to create document record", "error", err.Error())
        http.Error(w, "Failed to create document record: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Return response
    response := models.UploadResponse{
        UploadURL:  uploadURL,
        DocumentID: documentID,
        UserID: userID,
        Key:        key,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func UploadCompleteHandler(w http.ResponseWriter, r *http.Request) {
    var payload models.UploadPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        slog.Error("invalid request", "error", err)
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    slog.Info("saving document record")    
    err := services.SaveDocumentRecord(payload.DocumentID)
    if err != nil {
        slog.Error("failed to save document record", "error", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    slog.Info("sending payload to python worker")
    err = notifyWorker(payload)
    if err != nil {
        slog.Error("failed to send to queue", "error", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}


func notifyWorker(payload models.UploadPayload) error {
    workerURL := "http://worker:8080/process-document"

    body, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", workerURL, bytes.NewBuffer(body))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("worker returned status: %s", resp.Status)
    }

    return nil
}