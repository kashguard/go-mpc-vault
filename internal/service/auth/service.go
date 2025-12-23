package auth

import (
	"context"
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type AuthService interface {
	BeginRegistration(ctx context.Context, user string) (*protocol.CredentialCreation, *webauthn.SessionData, error)
	FinishRegistration(ctx context.Context, user string, sessionData webauthn.SessionData, response *http.Request) (*webauthn.Credential, error)
	BeginLogin(ctx context.Context, user string) (*protocol.CredentialAssertion, *webauthn.SessionData, error)
	FinishLogin(ctx context.Context, user string, sessionData webauthn.SessionData, response *http.Request) (*webauthn.Credential, error)
}
