package vault

import (
	"context"

	"github.com/campbell-frost/mailfs/internal/index"
)

func (v *Vault) Status(ctx context.Context, id string) (index.FileStatus, error) {
	return v.idx.File.Status(ctx, id)
}
