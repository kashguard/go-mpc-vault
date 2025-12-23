package api

import (
	"database/sql"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/wire"
	"google.golang.org/grpc"

	"github.com/kashguard/go-mpc-vault/internal/api/grpc/server"
	"github.com/kashguard/go-mpc-vault/internal/config"
	"github.com/kashguard/go-mpc-vault/internal/infra/mpc"
	mpcAuth "github.com/kashguard/go-mpc-vault/internal/service/auth"
	"github.com/kashguard/go-mpc-vault/internal/service/signing"
	"github.com/kashguard/go-mpc-vault/internal/service/vault"
)

var MpcProviderSet = wire.NewSet(
	NewMpcClientConnection,
	NewKeyClient,
	NewSigningClient,
	NewMpcAuthService,
	NewVaultService,
	NewSigningService,
	NewGrpcServer,
)

func NewMpcClientConnection(cfg config.Server) (*grpc.ClientConn, error) {
	if cfg.Mpc.CACertFile == "" {
		// In test environment or when MPC is not configured, return nil connection.
		// Services depending on this should handle nil connection or not be invoked.
		return nil, nil
	}

	mpcCfg := mpc.Config{
		Address:    cfg.Mpc.Address,
		ServerName: cfg.Mpc.ServerName,
		CACertFile: cfg.Mpc.CACertFile,
		CertFile:   cfg.Mpc.CertFile,
		KeyFile:    cfg.Mpc.KeyFile,
	}
	return mpc.NewClientConnection(mpcCfg)
}

func NewKeyClient(conn *grpc.ClientConn) *mpc.KeyClient {
	return mpc.NewKeyClient(conn)
}

func NewSigningClient(conn *grpc.ClientConn) *mpc.SigningClient {
	return mpc.NewSigningClient(conn)
}

func NewMpcAuthService(cfg config.Server, db *sql.DB) (mpcAuth.AuthService, error) {
	// TODO: Move these to config
	wconfig := &webauthn.Config{
		RPDisplayName: "MPC Vault",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:3000"},
	}
	w, err := webauthn.New(wconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to init webauthn: %w", err)
	}
	return mpcAuth.NewService(db, w), nil
}

//nolint:ireturn
func NewVaultService(db *sql.DB, keyClient *mpc.KeyClient) vault.Service {
	return vault.NewService(db, keyClient)
}

//nolint:ireturn
func NewSigningService(db *sql.DB, signingClient *mpc.SigningClient) signing.Service {
	return signing.NewService(db, signingClient)
}

func NewGrpcServer(
	cfg config.Server,
	authSvc mpcAuth.AuthService,
	vaultSvc vault.Service,
	signingSvc signing.Service,
) *grpc.Server {
	s := grpc.NewServer()

	authServer := server.NewAuthServer(authSvc)
	server.RegisterAuthServer(s, authServer)

	vaultServer := server.NewVaultServer(vaultSvc)
	server.RegisterVaultServer(s, vaultServer)

	signingServer := server.NewSigningServer(signingSvc)
	server.RegisterSigningServer(s, signingServer)

	return s
}
