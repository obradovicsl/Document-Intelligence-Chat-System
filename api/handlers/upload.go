package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/auth"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/dto"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/services"
)

type DocumentHandler struct {
	service *services.DocumentService
}

func NewDocumentHandler(service *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		service: service,
	}
}

// POST /api/upload/init
func (h *DocumentHandler) HandleInitUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(r)

	slog.Info("init upload request", "user_id", userID)

	// Decode request body
	var req dto.UploadDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("invalid request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FileName == "" {
		slog.Error("file name is required")
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	if req.FileSize <= 0 {
		slog.Error("file size must be positive")
		http.Error(w, "Invalid file size", http.StatusBadRequest)
		return
	}

	// Create document & presignedURL
	doc, uploadURL, err := h.service.CreateDocument(
		ctx,
		userID,
		req.FileName,
		req.FileType,
		req.FileSize,
	)

	if err != nil {
		slog.Error("failed to create document", "error", err)
		http.Error(w, "Failed to initialize upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Vrati response
	response := dto.UploadDocumentResponse{
		UploadURL:  uploadURL,
		DocumentID: doc.ID,
		UserID:     doc.UserID,
		Key:        doc.S3Key,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	slog.Info("upload initialized",
		"document_id", doc.ID,
		"user_id", userID,
		"file_name", req.FileName)
}

// POST /api/upload/complete
func (h *DocumentHandler) HandleCompleteUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(r)

	slog.Info("upload complete request", "user_id", userID)

	// Decode payload
	var payload dto.UploadDocumentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Error("invalid request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	doc, err := h.service.GetDocumentByID(ctx, payload.DocumentID)
	if err != nil {
		slog.Error("failed to get document", "error", err, "document_id", payload.DocumentID)
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	if doc.UserID != userID {
		slog.Warn("unauthorized upload complete attempt",
			"document_id", payload.DocumentID,
			"owner_id", doc.UserID,
			"requester_id", userID)
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	slog.Info("marking document as uploaded", "document_id", payload.DocumentID)
	if err := h.service.MarkAsUploaded(ctx, payload.DocumentID); err != nil {
		slog.Error("failed to mark as uploaded", "error", err, "document_id", payload.DocumentID)
		http.Error(w, "Failed to update document status", http.StatusInternalServerError)
		return
	}

	// Pošalji na Python worker za processing
	slog.Info("notifying worker", "document_id", payload.DocumentID)
	if err := h.notifyWorker(payload); err != nil {
		slog.Error("failed to notify worker", "error", err, "document_id", payload.DocumentID)
		
		http.Error(w, "Failed to start processing", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "processing",
		"message": "Document processing started",
	})

	slog.Info("upload completed successfully",
		"document_id", payload.DocumentID,
		"user_id", userID)
}

// notifyWorker šalje payload Python workeru za processing
func (h *DocumentHandler) notifyWorker(payload dto.UploadDocumentPayload) error {
	workerURL := "http://worker:8080/documents/process"

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", workerURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to worker: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("worker returned status: %s", resp.Status)
	}

	return nil
}