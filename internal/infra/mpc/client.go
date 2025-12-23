package mpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	Address    string
	ServerName string
	CACertFile string
	CertFile   string
	KeyFile    string
}

func NewClientConnection(cfg Config) (*grpc.ClientConn, error) {
	// Load CA cert
	caCert, err := os.ReadFile(cfg.CACertFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CA cert")
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("failed to append CA cert")
	}

	// Load client cert/key
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load client cert/key")
	}

	// Create TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ServerName:   cfg.ServerName,
		MinVersion:   tls.VersionTLS12,
	})

	//nolint:staticcheck // grpc.Dial is deprecated but used here for simplicity in PoC
	conn, err := grpc.Dial(cfg.Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MPC server: %w", err)
	}

	return conn, nil
}
