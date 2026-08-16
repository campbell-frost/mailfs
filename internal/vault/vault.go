package vault

import (
	"github.com/campbell-frost/mailfs/internal/gmail"
	"github.com/campbell-frost/mailfs/internal/index"
)

type Vault struct {
	idx     *index.Index
	gmail   *gmail.Gmail
	tempDir string
}

func New(
	index *index.Index,
	gmail *gmail.Gmail,
	tempDir string) *Vault {
	return &Vault{
		idx:     index,
		gmail:   gmail,
		tempDir: tempDir,
	}
}
