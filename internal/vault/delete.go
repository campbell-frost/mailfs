package vault

import (
	"context"
	"log"
)

func (v *Vault) Delete(ctx context.Context, id string) error {
	chunks, err := v.idx.Chunk.Get(ctx, id)
	if err != nil {
		return err
	}

	chunkRefs := make([][]byte, 0, len(chunks))
	for _, c := range chunks {
		chunkRefs = append(chunkRefs, c.Ref)
	}
	err = v.gmail.Delete(ctx, chunkRefs)
	if err != nil {
		log.Println("failed to delete chunks from gmail", err.Error())
		return err
	}
	log.Printf("deleted %d chunks for %s", len(chunks), id)

	if err := v.idx.File.Delete(ctx, id); err != nil {
		log.Println("failed to delete index", err.Error())
		return err
	}
	return nil
}
