package vault

import (
	"context"

	"github.com/kashguard/go-mpc-vault/internal/models"
)

type Service interface {
	CreateVault(ctx context.Context, name string, orgID string) (*models.Vault, error)
	CreateWallet(ctx context.Context, vaultID string, chainID string) (*models.Wallet, error)
}
