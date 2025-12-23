package signing

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/google/uuid"
	"github.com/kashguard/go-mpc-vault/internal/infra/mpc"
	"github.com/kashguard/go-mpc-vault/internal/models"
)

type impl struct {
	db            *sql.DB
	signingClient *mpc.SigningClient
}

//nolint:ireturn
func NewService(db *sql.DB, signingClient *mpc.SigningClient) Service {
	return &impl{
		db:            db,
		signingClient: signingClient,
	}
}

func (s *impl) CreateRequest(ctx context.Context, walletID string, txData string, note string, userID string) (*models.SigningRequest, error) {
	req := &models.SigningRequest{
		ID:          uuid.New().String(),
		WalletID:    null.StringFrom(walletID),
		TXData:      txData,
		Note:        null.StringFrom(note),
		Status:      null.StringFrom("pending"),
		InitiatorID: null.StringFrom(userID),
	}

	if err := req.Insert(ctx, s.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("failed to insert signing request: %w", err)
	}

	return req, nil
}

func (s *impl) ApproveRequest(ctx context.Context, requestID string, params ApprovalParams) error {
	// 1. Load Request with Wallet and Vault
	req, err := models.SigningRequests(
		models.SigningRequestWhere.ID.EQ(requestID),
		qm.Load(models.SigningRequestRels.Wallet),
		qm.Load(qm.Rels(models.SigningRequestRels.Wallet, models.WalletRels.Vault)),
		qm.Load(qm.Rels(models.SigningRequestRels.Wallet, models.WalletRels.Chain)),
	).One(ctx, s.db)
	if err != nil {
		return fmt.Errorf("request not found: %w", err)
	}

	if req.Status.String != "pending" {
		return fmt.Errorf("request is not pending (status: %s)", req.Status.String)
	}

	// 2. Check if user already approved
	exists, err := models.Approvals(
		models.ApprovalWhere.RequestID.EQ(null.StringFrom(requestID)),
		models.ApprovalWhere.UserID.EQ(null.StringFrom(params.UserID)),
	).Exists(ctx, s.db)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user already approved")
	}

	// 3. Store Approval
	authDataJSON, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}
	approval := &models.Approval{
		ID:        uuid.New().String(),
		RequestID: null.StringFrom(requestID),
		UserID:    null.StringFrom(params.UserID),
		Action:    "approve",
		Comment:   null.StringFrom(string(authDataJSON)), // Storing auth data in comment
	}
	if err := approval.Insert(ctx, s.db, boil.Infer()); err != nil {
		return err
	}

	// 4. Check Threshold
	approvals, err := models.Approvals(
		models.ApprovalWhere.RequestID.EQ(null.StringFrom(requestID)),
		models.ApprovalWhere.Action.EQ("approve"),
	).All(ctx, s.db)
	if err != nil {
		return err
	}

	vault := req.R.Wallet.R.Vault
	if vault == nil {
		return fmt.Errorf("vault not found for wallet")
	}

	if len(approvals) >= vault.Threshold {
		// Trigger MPC
		var authTokens []mpc.AuthToken
		for _, app := range approvals {
			var param ApprovalParams
			if app.Comment.Valid {
				_ = json.Unmarshal([]byte(app.Comment.String), &param)
				authTokens = append(authTokens, mpc.AuthToken{
					PasskeySignature:  param.Signature,
					AuthenticatorData: param.AuthenticatorData,
					ClientDataJson:    param.ClientDataJSON,
					CredentialId:      param.CredentialID,
				})
			}
		}

		wallet := req.R.Wallet
		chain := wallet.R.Chain
		chainType := "evm" // Default
		if chain != nil {
			chainType = chain.Type
		}

		signature, err := s.signingClient.ThresholdSign(ctx, wallet.KeyID, req.TXData, chainType, authTokens)
		if err != nil {
			// Mark as failed? Or just return error?
			// Ideally we should log this and maybe retry or mark as failed.
			return fmt.Errorf("mpc signing failed: %w", err)
		}

		// Update Request
		req.Status = null.StringFrom("completed")
		req.Signature = null.StringFrom(signature)
		if _, err := req.Update(ctx, s.db, boil.Infer()); err != nil {
			return err
		}
	}

	return nil
}

func (s *impl) RejectRequest(ctx context.Context, requestID string, userID string) error {
	req, err := models.FindSigningRequest(ctx, s.db, requestID)
	if err != nil {
		return err
	}

	if req.Status.String != "pending" {
		return fmt.Errorf("request is not pending")
	}

	approval := &models.Approval{
		ID:        uuid.New().String(),
		RequestID: null.StringFrom(requestID),
		UserID:    null.StringFrom(userID),
		Action:    "reject",
	}
	if err := approval.Insert(ctx, s.db, boil.Infer()); err != nil {
		return err
	}

	req.Status = null.StringFrom("rejected")
	if _, err := req.Update(ctx, s.db, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func (s *impl) GetRequest(ctx context.Context, requestID string) (*models.SigningRequest, error) {
	return models.FindSigningRequest(ctx, s.db, requestID)
}

func (s *impl) ListRequests(ctx context.Context, userID string, vaultID string, status string, page int, limit int) (models.SigningRequestSlice, int64, error) {
	mods := []qm.QueryMod{
		qm.InnerJoin("vaults on vaults.id = signing_requests.vault_id"),
		qm.InnerJoin("organization_members on organization_members.organization_id = vaults.organization_id"),
		qm.Where("organization_members.user_id = ?", userID),
	}
	if status != "" {
		mods = append(mods, models.SigningRequestWhere.Status.EQ(null.StringFrom(status)))
	}
	if vaultID != "" {
		mods = append(mods, models.SigningRequestWhere.VaultID.EQ(null.StringFrom(vaultID)))
	}
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	total, err := models.SigningRequests(mods...).Count(ctx, s.db)
	if err != nil {
		return nil, 0, err
	}

	mods = append(mods, qm.OrderBy("signing_requests.created_at DESC"), qm.Limit(limit), qm.Offset(offset))
	items, err := models.SigningRequests(mods...).All(ctx, s.db)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
