package api

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
)

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "form too large", http.StatusBadRequest)
		return
	}

	defer r.MultipartForm.RemoveAll()

	f, h, err := r.FormFile("mailfs.file")
	if err != nil {
		http.Error(w, "missing file in form", http.StatusBadRequest)
		return
	}
	defer f.Close()

	id, err := s.vault.Stage(r.Context(), filepath.Base(h.Filename), f)
	if err != nil {
		http.Error(w, "failed to stage file", http.StatusInternalServerError)
		return
	}

	log.Printf("created temp file for %s, id: %s\n", h.Filename, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}
