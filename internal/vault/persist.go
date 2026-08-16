package vault

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/campbell-frost/mailfs/internal/index"
)

// Persist takes a staged file, chunks it according to the configured chunk limit,
// and stores the data in imap messages, and index data in the sqlite chunks table
func (v *Vault) Persist(ctx context.Context, id string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("persist panicked: %v", r)
		}
		if err != nil {
			log.Println("persist failed", err)
			v.idx.SetStatus(context.WithoutCancel(ctx), id, index.StatusFailed)
		}
	}()

	if err := v.idx.SetStatus(ctx, id, index.StatusProcessing); err != nil {
		return err
	}

	fi, err := v.idx.Stat(ctx, id)
	if err != nil {
		return err
	}

	tempPath := filepath.Join(fi.TempPath, fi.ID)

	f, err := os.Open(tempPath)
	if err != nil {
		return err
	}

	defer f.Close()

	chunkSize := int64(v.gmail.MaxChunkSize())
	count := int((fi.Size + chunkSize - 1) / chunkSize)

	if count == 0 {
		return errors.New("no chunks to persist")
	}

	for seq := range count {
		if err := v.persistChunk(ctx, f, fi, seq, chunkSize); err != nil {
			return fmt.Errorf("chunk %d: %w", seq, err)
		}
	}

	if err := v.idx.SetStatus(ctx, id, index.StatusStored); err != nil {
		return err
	}
	if err := os.Remove(tempPath); err != nil {
		log.Println("temp file not removed", err)
	}
	return nil
}

func (v *Vault) persistChunk(ctx context.Context, f *os.File, fi index.FileInfo, seq int, chunkSize int64) error {
	off := int64(seq) * chunkSize
	data := make([]byte, min(chunkSize, fi.Size-off))
	if _, err := f.ReadAt(data, off); err != nil {
		return fmt.Errorf("failed to read chunk at %d: %w", off, err)
	}

	sum := sha256.Sum256(data)
	_ = hex.EncodeToString(sum[:])

	// put message to mailbox
	// persist to index

	return nil
}
