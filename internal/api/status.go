package api

import (
	"encoding/json"
	"net/http"
)

type statusResponse struct {
	Status string `json:"status"`
	// ProcessedChunks int `json:"processedChunks"`
}

func (s *Server) StatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	st, err := s.vault.Status(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to get status", http.StatusBadRequest)
	}

	resp := statusResponse{Status: st}
	json.NewEncoder(w).Encode(resp)
}
