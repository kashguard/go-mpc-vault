package signing

import (
	"context"

	"github.com/kashguard/go-mpc-vault/internal/models"
)

type ApprovalParams struct {
	UserID            string
	CredentialID      []byte
	Signature         []byte
	AuthenticatorData []byte
	ClientDataJSON    []byte
}

type Service interface {
	CreateRequest(ctx context.Context, walletID string, txData string, note string, userID string) (*models.SigningRequest, error)
	ApproveRequest(ctx context.Context, requestID string, params ApprovalParams) error
	RejectRequest(ctx context.Context, requestID string, userID string) error
	GetRequest(ctx context.Context, requestID string) (*models.SigningRequest, error)
	ListRequests(ctx context.Context, userID string, vaultID string, status string, page int, limit int) (models.SigningRequestSlice, int64, error)
}
