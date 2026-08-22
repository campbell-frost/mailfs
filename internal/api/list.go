package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) ListHandler(w http.ResponseWriter, r *http.Request) {
	files, err := s.vault.ListFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := make([]fileInfoResponse, 0, len(files))
	for _, file := range files {
		out = append(out, fileInfoResponse{
			ID:        file.ID,
			Filename:  file.Filename,
			Size:      file.Size,
			Status:    file.Status,
			CreatedAt: file.CreatedAt,
			Chunks:    file.Chunks,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(out)
}
