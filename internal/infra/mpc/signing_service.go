package mpc

import (
	"context"
	"encoding/base64"

	infra "github.com/kashguard/go-mpc-vault/internal/infra/grpc/v1"
	"google.golang.org/grpc"
)

type SigningClient struct {
	client infra.SigningServiceClient
}

func NewSigningClient(conn *grpc.ClientConn) *SigningClient {
	return &SigningClient{
		client: infra.NewSigningServiceClient(conn),
	}
}

type AuthToken struct {
	PasskeySignature  []byte
	AuthenticatorData []byte
	ClientDataJson    []byte
	CredentialId      []byte
}

func (c *SigningClient) ThresholdSign(ctx context.Context, keyID string, messageHex string, chainType string, authTokens []AuthToken) (string, error) {
	infraTokens := make([]*infra.AuthToken, len(authTokens))
	for i, t := range authTokens {
		infraTokens[i] = &infra.AuthToken{
			PasskeySignature:  t.PasskeySignature,
			AuthenticatorData: t.AuthenticatorData,
			ClientDataJson:    t.ClientDataJson,
			CredentialId:      base64.RawURLEncoding.EncodeToString(t.CredentialId),
		}
	}

	req := &infra.ThresholdSignRequest{
		KeyId:      keyID,
		MessageHex: messageHex,
		ChainType:  chainType,
		AuthTokens: infraTokens,
	}

	resp, err := c.client.ThresholdSign(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetSignature(), nil
}
