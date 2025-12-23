package vault

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/google/uuid"
	"github.com/kashguard/go-mpc-vault/internal/infra/mpc"
	"github.com/kashguard/go-mpc-vault/internal/models"
)

type impl struct {
	db        *sql.DB
	keyClient *mpc.KeyClient
}

//nolint:ireturn
func NewService(db *sql.DB, keyClient *mpc.KeyClient) Service {
	return &impl{
		db:        db,
		keyClient: keyClient,
	}
}

func (s *impl) CreateVault(ctx context.Context, name string, orgID string) (*models.Vault, error) {
	vault := &models.Vault{
		Name:           name,
		OrganizationID: null.StringFrom(orgID),
		Threshold:      2, // Default
	}

	if err := vault.Insert(ctx, s.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to insert vault: %w", err)
	}

	return vault, nil
}

func (s *impl) CreateWallet(ctx context.Context, vaultID string, chainID string) (*models.Wallet, error) {
	// 1. Check if chain exists
	chain, err := models.Chains(models.ChainWhere.ID.EQ(chainID)).One(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("chain not found: %w", err)
	}

	// 2. Generate ID
	walletID := uuid.New().String()

	// 3. Call MPC to generate key
	// Use wallet ID as Key ID
	// Determine algo/curve based on chain
	algo := "ECDSA"
	curve := "secp256k1"
	if chain.Curve == "ed25519" { // Assuming Chain has Curve field or similar
		curve = "ed25519"
		algo = "EdDSA"
	}

	pubKey, err := s.keyClient.CreateKey(ctx, walletID, algo, curve)
	if err != nil {
		return nil, fmt.Errorf("mpc key generation failed: %w", err)
	}

	// 4. Create Wallet record
	wallet := &models.Wallet{
		ID:          walletID,
		VaultID:     null.StringFrom(vaultID),
		ChainID:     null.StringFrom(chainID),
		KeyID:       walletID, // Using walletID as KeyID
		Address:     deriveAddress(pubKey, chain.Type),
		DerivePath:  "m/44'/60'/0'/0/0", // Simplified default
		DeriveIndex: 0,
	}

	if err := wallet.Insert(ctx, s.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to insert wallet: %w", err)
	}

	return wallet, nil
}

func deriveAddress(pubKeyHex string, chainType string) string {
	// Placeholder for address derivation
	// Real impl needs crypto libs
	return "0x" + pubKeyHex[:40] // Mock
}
