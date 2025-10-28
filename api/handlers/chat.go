package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/auth"
	"github.com/obradovicsl/Document-Intelligence-Chat-System/API/dto"
)

// GET /api/documents
func (h *DocumentHandler) HandleQuestion(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)

	var req dto.ChatRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.UserID = userID

	slog.Debug("chat question received", "user_id", req.UserID, "question", req.Question)

	payload, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to serialize request", http.StatusInternalServerError)
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(h.pythonWorkerURL+"/chat/ask", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		slog.Error("failed to contact python worker", "error", err)
		http.Error(w, "Failed to contact worker", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("python worker error", "status", resp.StatusCode, "body", string(body))
		http.Error(w, "Worker returned error", http.StatusBadGateway)
		return
	}

	var chatResp dto.ChatResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		http.Error(w, "Failed to parse worker response", http.StatusInternalServerError)
		return
	}

	// --- Saljemo odgovor frontu ---
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResp)
}