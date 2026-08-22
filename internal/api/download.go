package api

import (
	"net/http"

	"github.com/campbell-frost/mailfs/internal/index"
)

func (s *Server) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fi, err := s.vault.Lookup(r.Context(), id)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	if fi.Status != index.StatusStored {
		http.Error(w, "invalid file status", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	if err := s.vault.Fetch(r.Context(), id, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
