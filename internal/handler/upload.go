package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const TEMP_DIR = "tmp"

func UploadHandler(w http.ResponseWriter, r *http.Request) {
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

	// create temp file on disk
	if err := os.MkdirAll(TEMP_DIR, 0o755); err != nil {
		http.Error(w, "failed to create temp file", http.StatusInternalServerError)
		return
	}

	id := uuid.NewString()
	tempPath := filepath.Join(TEMP_DIR, id)

	tmp, err := os.Create(tempPath)
	if err != nil {
		http.Error(w, "failed to create temp file", http.StatusInternalServerError)
		return
	}

	defer tmp.Close()
	n, err := io.Copy(tmp, f)
	if err != nil {
		http.Error(w, "failed to create temp file", http.StatusInternalServerError)
		return
	}
	fmt.Printf("created temp file for %s with size %d\n", h.Filename, n)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}
