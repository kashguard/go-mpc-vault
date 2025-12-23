package api

import (
	"database/sql"

	"github.com/kashguard/go-mpc-vault/internal/service/organization"
)

func NewOrganizationService(db *sql.DB) organization.Service {
	return organization.NewService(db)
}
