package mpc

import (
	"context"

	infra "github.com/kashguard/go-mpc-vault/internal/infra/grpc/v1"
	"google.golang.org/grpc"
)

type KeyClient struct {
	client infra.KeyServiceClient
}

func NewKeyClient(conn *grpc.ClientConn) *KeyClient {
	return &KeyClient{
		client: infra.NewKeyServiceClient(conn),
	}
}

func (c *KeyClient) CreateKey(ctx context.Context, keyID string, algorithm string, curve string) (string, error) {
	req := &infra.CreateRootKeyRequest{
		KeyId:     keyID,
		Algorithm: algorithm,
		Curve:     curve,
	}

	resp, err := c.client.CreateRootKey(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetKey().GetPublicKey(), nil
}
