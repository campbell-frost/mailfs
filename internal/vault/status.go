package vault

import "context"

func (v *Vault) Status(ctx context.Context, id string) (string, error) {
	return v.idx.Status(ctx, id)
}
