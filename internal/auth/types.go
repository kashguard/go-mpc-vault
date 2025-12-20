package auth

import (
	"time"

	"github.com/kashguard/go-mpc-vault/internal/data/dto"
)

type Result struct {
	Token      string
	User       *dto.User
	ValidUntil time.Time
	Scopes     []string
}
