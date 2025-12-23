package server

import (
	"context"

	apiv1 "github.com/kashguard/go-mpc-vault/internal/api/grpc/v1"
	"github.com/kashguard/go-mpc-vault/internal/service/signing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SigningServer struct {
	apiv1.UnimplementedSigningServiceServer
	service signing.Service
}

func NewSigningServer(s signing.Service) *SigningServer {
	return &SigningServer{
		service: s,
	}
}

func (s *SigningServer) CreateSigningRequest(ctx context.Context, req *apiv1.CreateSigningRequest) (*apiv1.CreateSigningResponse, error) {
	// TODO: Get UserID from context
	userID := "default-user"

	r, err := s.service.CreateRequest(ctx, req.GetWalletId(), req.GetTxData(), req.GetNote(), userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create request: %v", err)
	}

	return &apiv1.CreateSigningResponse{
		RequestId: r.ID,
		Status:    r.Status.String,
	}, nil
}

func (s *SigningServer) ApproveSigningRequest(ctx context.Context, req *apiv1.ApproveSigningRequest) (*apiv1.ApproveSigningResponse, error) {
	// TODO: Get UserID from context
	userID := "default-user"

	params := signing.ApprovalParams{
		UserID:            userID,
		CredentialID:      req.GetCredentialId(),
		Signature:         req.GetSignature(),
		AuthenticatorData: req.GetAuthenticatorData(),
		ClientDataJSON:    req.GetClientDataJson(),
	}

	// Workaround: We'll assume the client sends a JSON string in `Comment` containing all auth data
	// if we can't update proto easily.
	// But let's check if we can update proto. We generated it earlier.
	// Wait, the `ApproveSigningRequest` message in `api/v1/signing.proto` has:
	// RequestId, Action, Comment.
	// It DOES NOT have CredentialId, Signature, AuthenticatorData, ClientDataJson.

	// I should update the proto file first.

	if err := s.service.ApproveRequest(ctx, req.GetRequestId(), params); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to approve request: %v", err)
	}

	return &apiv1.ApproveSigningResponse{
		Status: "approved",
	}, nil
}

func RegisterSigningServer(s *grpc.Server, srv *SigningServer) {
	apiv1.RegisterSigningServiceServer(s, srv)
}
