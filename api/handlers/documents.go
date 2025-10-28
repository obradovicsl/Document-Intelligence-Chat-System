package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/auth"
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"documents": docs,
		"count":     len(docs),
	})
}

// GET /api/documents/{id}
func (h *DocumentHandler) HandleGetDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := auth.GetUserID(r)
	documentID := r.PathValue("id") // Go 1.22+ path parameter

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

	// Proveri vlasništvo
	if doc.UserID != userID {
		slog.Warn("unauthorized document access",
			"document_id", documentID,
			"owner_id", doc.UserID,
			"requester_id", userID)
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}