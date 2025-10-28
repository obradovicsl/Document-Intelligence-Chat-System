package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/auth"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/dto"
)

// GET /api/documents
func (h *DocumentHandler) HandleGetDocuments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(r)

	slog.Debug("get documents request", "user_id", userID)

	docs, err := h.service.GetUserDocuments(ctx, userID)
	if err != nil {
		slog.Error("failed to get documents", "error", err, "user_id", userID)
		http.Error(w, "Failed to fetch documents", http.StatusInternalServerError)
		return
	}

	docsDTO := dto.ToDocumentDTOList(docs)

	response := dto.DocumentListResponse{
		Documents: docsDTO,
		Count: len(docsDTO),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/documents/{id}
func (h *DocumentHandler) HandleGetDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(r)
	documentID := r.PathValue("id")

	if documentID == "" {
		http.Error(w, "Document ID is required", http.StatusBadRequest)
		return
	}

	doc, err := h.service.GetDocumentByID(ctx, documentID)
	if err != nil {
		slog.Error("failed to get document", "error", err, "document_id", documentID)
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	if doc.UserID != userID {
		slog.Warn("unauthorized document access",
			"document_id", documentID,
			"owner_id", doc.UserID,
			"requester_id", userID)
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	response := dto.ToDocumentDTO(doc)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/documents/user/me
func (h *DocumentHandler) HandleGetDocumentForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUserID := auth.GetUserID(r)

	docs, err := h.service.GetUserDocuments(ctx, currentUserID)
	if err != nil {
		slog.Error("failed to get document", "error", err, "user_id", currentUserID)
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	docsDTO := dto.ToDocumentDTOList(docs)

	response := dto.DocumentListResponse{
		Documents: docsDTO,
		Count: len(docsDTO),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}