package server

import (
	"context"

	apiv1 "github.com/kashguard/go-mpc-vault/internal/api/grpc/v1"
	"github.com/kashguard/go-mpc-vault/internal/service/vault"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VaultServer struct {
	apiv1.UnimplementedVaultServiceServer
	service vault.Service
}

func NewVaultServer(s vault.Service) *VaultServer {
	return &VaultServer{
		service: s,
	}
}

func (s *VaultServer) CreateVault(ctx context.Context, req *apiv1.CreateVaultRequest) (*apiv1.CreateVaultResponse, error) {
	// TODO: Get OrgID from context
	orgID := "default-org"

	v, err := s.service.CreateVault(ctx, req.GetName(), orgID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create vault: %v", err)
	}

	// TODO: Handle threshold and chains initialization if service supports it
	// Current service implementation only inserts vault record.
	// We might need to extend service to support initial wallets.

	if len(req.GetChains()) > 0 {
		for _, chainID := range req.GetChains() {
			_, err := s.service.CreateWallet(ctx, v.ID, chainID)
			if err != nil {
				// Log error but don't fail entire request? Or fail?
				// For now, fail.
				return nil, status.Errorf(codes.Internal, "failed to create wallet for chain %s: %v", chainID, err)
			}
		}
	}

	return &apiv1.CreateVaultResponse{
		VaultId: v.ID,
		Status:  "active",
	}, nil
}

func (s *VaultServer) CreateWallet(ctx context.Context, req *apiv1.CreateWalletRequest) (*apiv1.CreateWalletResponse, error) {
	w, err := s.service.CreateWallet(ctx, req.GetVaultId(), req.GetChainId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create wallet: %v", err)
	}

	return &apiv1.CreateWalletResponse{
		WalletId: w.ID,
		Address:  w.Address,
		ChainId:  w.ChainID.String,
	}, nil
}

func RegisterVaultServer(s *grpc.Server, srv *VaultServer) {
	apiv1.RegisterVaultServiceServer(s, srv)
}
