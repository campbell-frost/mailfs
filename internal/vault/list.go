package vault

import "github.com/campbell-frost/mailfs/internal/index"

func (v *Vault) ListFiles() ([]index.FileInfo, error) {
	return v.idx.Files()
}
