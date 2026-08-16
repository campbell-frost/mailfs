package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "expected a multipart body", http.StatusBadRequest)
		return
	}

	part, err := findFilePart(mr, "mailfs.file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer part.Close()

	name := filepath.Base(part.FileName())
	id, err := s.vault.Stage(r.Context(), name, part)
	if err != nil {
		http.Error(w, "failed to stage file", http.StatusInternalServerError)
		return
	}

	log.Printf("created temp file for %s, id: %s\n", name, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func findFilePart(mr *multipart.Reader, name string) (*multipart.Part, error) {
	for {
		p, err := mr.NextPart()
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("no %v part in request", name)
		}
		if err != nil {
			return nil, err
		}

		if p.FormName() == name && p.FileName() != "" {
			return p, nil
		}
		p.Close()
	}
}
