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
	"time"

	"github.com/campbell-frost/mailfs/internal/gmail"
	"github.com/campbell-frost/mailfs/internal/index"
	"golang.org/x/sync/errgroup"
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

	count := v.chunkCount(fi.Size)

	if count == 0 {
		return errors.New("no chunks to persist")
	}

	workers := min(v.gmail.Concurrency(), count)
	start := time.Now()

	log.Printf("persist start file=%s name=%q size=%d chunks=%d workers=%d",
		fi.ID, fi.Filename, fi.Size, count, workers,
	)

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(workers)

	for seq := range count {
		g.Go(func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("chunk %d panicked: %v", seq, r)
				}
			}()

			chunkStart := time.Now()

			if err := v.persistChunk(gCtx, f, fi, seq); err != nil {
				return fmt.Errorf("chunk %d: %w", seq, err)
			}

			log.Printf("persist chunk file=%s seq=%d of %d bytes=%d took=%s",
				fi.ID, seq+1, count, v.gmail.ChunkSize(), time.Since(chunkStart).Round(time.Millisecond))
			return nil

		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	if err := v.idx.SetStatus(ctx, id, index.StatusStored); err != nil {
		return err
	}
	if err := os.Remove(tempPath); err != nil {
		log.Println("temp file not removed", err)
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	log.Printf("persist done  file=%s chunks=%d bytes=%d took=%s", fi.ID, count, fi.Size, elapsed)
	return nil
}

func (v *Vault) persistChunk(ctx context.Context, f *os.File, fi index.FileInfo, seq int) error {
	chunkSize := v.gmail.ChunkSize()
	off := int64(seq) * chunkSize
	data := make([]byte, min(chunkSize, fi.Size-off))
	if _, err := f.ReadAt(data, off); err != nil {
		return fmt.Errorf("failed to read chunk at %d: %w", off, err)
	}

	sum := sha256.Sum256(data)
	sha := hex.EncodeToString(sum[:])

	ref, err := v.gmail.Put(ctx, gmail.Chunk{
		FileID: fi.ID,
		Seq:    seq,
		Size:   len(data),
		Sha256: sha,
		Data:   data,
	})
	if err != nil {
		return fmt.Errorf("failed to store chunk via gmail: %w", err)
	}

	return v.idx.AddChunk(ctx, index.Chunk{
		FileID: fi.ID,
		Seq:    seq,
		Size:   len(data),
		Sha256: sha,
		Ref:    ref,
	})
}
