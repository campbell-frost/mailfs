package vault

import (
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
