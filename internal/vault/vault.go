package vault

import (
	"context"

	"github.com/campbell-frost/mailfs/internal/gmail"
	"github.com/campbell-frost/mailfs/internal/index"
)

type Vault struct {
	idx     *index.Index
	gmail   *gmail.Client
	tempDir string
}

func New(
	index *index.Index,
	client *gmail.Client,
	tempDir string) *Vault {
	return &Vault{
		idx:     index,
		gmail:   client,
		tempDir: tempDir,
	}
}

func (v *Vault) ListFiles() ([]index.FileInfo, error) {
	return v.idx.Files()
}

func (v *Vault) Lookup(ctx context.Context, id string) (index.FileInfo, error) {
	return v.idx.Stat(ctx, id)
}
