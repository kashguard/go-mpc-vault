package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/kashguard/go-mpc-vault/internal/models"
)

type Service struct {
	db       *sql.DB
	webAuthn *webauthn.WebAuthn
}

func NewService(db *sql.DB, w *webauthn.WebAuthn) *Service {
	return &Service{
		db:       db,
		webAuthn: w,
	}
}

// User wrapper to implement webauthn.User
type WebAuthnUser struct {
	*models.User
	credentials []webauthn.Credential
}

func (u *WebAuthnUser) WebAuthnID() []byte {
	return []byte(u.ID)
}

func (u *WebAuthnUser) WebAuthnName() string {
	return u.Username.String
}

func (u *WebAuthnUser) WebAuthnDisplayName() string {
	return u.Username.String
}

func (u *WebAuthnUser) WebAuthnIcon() string {
	return ""
}

func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func (s *Service) getUser(ctx context.Context, email string) (*WebAuthnUser, error) {
	user, err := models.Users(models.UserWhere.Username.EQ(null.StringFrom(email))).One(ctx, s.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Load credentials
	creds, err := models.UserCredentials(models.UserCredentialWhere.UserID.EQ(user.ID)).All(ctx, s.db)
	if err != nil {
		return nil, err
	}

	var webAuthnCreds []webauthn.Credential
	for _, c := range creds {
		// In a real implementation, we need to store and retrieve the full Credential struct or sufficient parts
		// For now, let's assume we decode from DB or construct it.
		// Since we only store PublicKey and ID in the migration, we might need to adjust the migration or reconstruction.
		// NOTE: The migration I created has `credential_id` and `public_key`.
		// `webauthn.Credential` has ID, PublicKey, AttestationType, etc.

		// For simplicity in this PoC, we might need to store the blob.
		// But let's try to reconstruct what we can.

		webAuthnCreds = append(webAuthnCreds, webauthn.Credential{
			ID:              []byte(c.CredentialID),
			PublicKey:       []byte(c.PublicKey),
			AttestationType: c.AttestationType.String,
			// Authenticator: ...
		})
	}

	return &WebAuthnUser{
		User:        user,
		credentials: webAuthnCreds,
	}, nil
}

func (s *Service) BeginRegistration(ctx context.Context, email string) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	user, err := s.getUser(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	return s.webAuthn.BeginRegistration(user)
}

func (s *Service) FinishRegistration(ctx context.Context, email string, sessionData webauthn.SessionData, r *http.Request) (*webauthn.Credential, error) {
	user, err := s.getUser(ctx, email)
	if err != nil {
		return nil, err
	}

	credential, err := s.webAuthn.FinishRegistration(user, sessionData, r)
	if err != nil {
		return nil, err
	}

	// Save credential to DB
	dbCred := &models.UserCredential{
		UserID:          user.ID,
		CredentialID:    string(credential.ID),
		PublicKey:       string(credential.PublicKey),
		AttestationType: null.StringFrom(credential.AttestationType),
		// SignCount:       int(credential.Authenticator.SignCount),
	}
	// Note: Authenticator AAGUID and SignCount are inside credential.Authenticator which is not exported fully in older versions?
	// Checking recent webauthn lib, Credential struct has Authenticator field.
	dbCred.Aaguid = null.StringFrom(string(credential.Authenticator.AAGUID))
	dbCred.SignCount = null.IntFrom(int(credential.Authenticator.SignCount))

	err = dbCred.Insert(ctx, s.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return credential, nil
}

func (s *Service) BeginLogin(ctx context.Context, email string) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	user, err := s.getUser(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	return s.webAuthn.BeginLogin(user)
}

func (s *Service) FinishLogin(ctx context.Context, email string, sessionData webauthn.SessionData, r *http.Request) (*webauthn.Credential, error) {
	user, err := s.getUser(ctx, email)
	if err != nil {
		return nil, err
	}

	credential, err := s.webAuthn.FinishLogin(user, sessionData, r)
	if err != nil {
		return nil, err
	}

	// Update sign count
	// We need to find the specific credential and update it
	dbCred, err := models.UserCredentials(
		models.UserCredentialWhere.UserID.EQ(user.ID),
		models.UserCredentialWhere.CredentialID.EQ(string(credential.ID)),
	).One(ctx, s.db)

	if err == nil {
		dbCred.SignCount = null.IntFrom(int(credential.Authenticator.SignCount))
		if _, err := dbCred.Update(ctx, s.db, boil.Infer()); err != nil {
			return nil, err
		}
	}

	return credential, nil
}
