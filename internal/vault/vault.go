package vault

import "github.com/campbell-frost/mailfs/internal/index"

type Vault struct {
	index   index.Index
	tempDir string
}

func New(index *index.Index, tempDir string) *Vault {
	return &Vault{
		index:   *index,
		tempDir: tempDir,
	}
}
