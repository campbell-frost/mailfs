package vault

import "github.com/campbell-frost/mailfs/internal/index"

type Vault struct {
	idx     index.Index
	tempDir string
}

func New(index *index.Index, tempDir string) *Vault {
	return &Vault{
		idx:     *index,
		tempDir: tempDir,
	}
}
