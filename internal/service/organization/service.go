package organization

import (
	"context"

	"github.com/kashguard/go-mpc-vault/internal/models"
)

type Service interface {
	CreateOrganization(ctx context.Context, name string, ownerID string) (*models.Organization, error)
	ListUserOrganizations(ctx context.Context, userID string) (models.OrganizationSlice, error)
	ListMembers(ctx context.Context, orgID string) (models.OrganizationMemberSlice, error)
	AddMember(ctx context.Context, orgID string, userID string, role string) (*models.OrganizationMember, error)
	RemoveMember(ctx context.Context, orgID string, userID string) error
	UpdateMemberRole(ctx context.Context, orgID string, userID string, role string) error
}
