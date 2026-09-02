package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type statusResponse struct {
	Status          string `json:"status"`
	ProcessedChunks int    `json:"processedChunks"`
}

func (s *Server) StatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	fs, err := s.vault.Status(r.Context(), id)
	if err != nil {
		log.Printf("status failed to get status %s, %e", id, err)
		http.Error(w, "failed to get status", http.StatusBadRequest)
		return
	}
	log.Printf("status id: %s status: %s processedChunks %d", id, fs.Status, fs.ProcessedChunks)

	resp := statusResponse{
		Status:          fs.Status,
		ProcessedChunks: fs.ProcessedChunks,
	}
	json.NewEncoder(w).Encode(resp)
}
