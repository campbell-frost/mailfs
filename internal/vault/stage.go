package vault

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/campbell-frost/mailfs/internal/index"
	"github.com/google/uuid"
)

// Stage creates a temp file on disk and creates a record in the sqlite db
func (v *Vault) Stage(ctx context.Context, name string, r io.Reader) (string, error) {
	if err := os.MkdirAll(v.tempDir, 0o755); err != nil {
		return "", err
	}

	id := uuid.NewString()
	tempPath := filepath.Join(v.tempDir, id)

	tmp, err := os.Create(tempPath)
	if err != nil {
		return "", err
	}

	defer tmp.Close()

	size, err := io.Copy(tmp, r)
	if err != nil {
		return "", err
	}

	f := index.FileInfo{
		ID:        id,
		Filename:  name,
		Size:      size,
		CreatedAt: time.Now(),
		Status:    index.StatusPending,
		TempPath:  v.tempDir,
	}

	err = v.idx.CreateFileInfo(ctx, f)
	if err != nil {
		return "", err
	}
	return id, nil
}
