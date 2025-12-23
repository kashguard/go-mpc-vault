package organization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/kashguard/go-mpc-vault/internal/models"
)

type impl struct {
	db *sql.DB
}

func NewService(db *sql.DB) Service {
	return &impl{
		db: db,
	}
}

func (s *impl) CreateOrganization(ctx context.Context, name string, ownerID string) (*models.Organization, error) {
	org := &models.Organization{
		Name:    name,
		OwnerID: ownerID,
	}

	if err := org.Insert(ctx, s.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert organization: %w", err)
	}
	return org, nil
}

func (s *impl) ListUserOrganizations(ctx context.Context, userID string) (models.OrganizationSlice, error) {
	owned, err := models.Organizations(models.OrganizationWhere.OwnerID.EQ(userID)).All(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("list owned organizations: %w", err)
	}

	memberOrgIDs, err := models.OrganizationMembers(models.OrganizationMemberWhere.UserID.EQ(userID)).All(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("list memberships: %w", err)
	}

	var ids []string
	for _, m := range memberOrgIDs {
		ids = append(ids, m.OrganizationID)
	}
	var memberOrgs models.OrganizationSlice
	if len(ids) > 0 {
		memberOrgs, err = models.Organizations(qm.WhereIn(models.OrganizationColumns.ID, ids)).All(ctx, s.db)
		if err != nil {
			return nil, fmt.Errorf("list member organizations: %w", err)
		}
	}

	return append(owned, memberOrgs...), nil
}

func (s *impl) ListMembers(ctx context.Context, orgID string) (models.OrganizationMemberSlice, error) {
	members, err := models.OrganizationMembers(models.OrganizationMemberWhere.OrganizationID.EQ(orgID)).All(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	return members, nil
}

func (s *impl) AddMember(ctx context.Context, orgID string, userID string, role string) (*models.OrganizationMember, error) {
	_, err := models.FindOrganization(ctx, s.db, orgID)
	if err != nil {
		return nil, fmt.Errorf("organization not found: %w", err)
	}

	member := &models.OrganizationMember{
		OrganizationID: orgID,
		UserID:         userID,
		Role:           role,
		CreatedAt:      null.Time{},
	}
	if err := member.Insert(ctx, s.db, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert member: %w", err)
	}
	return member, nil
}

func (s *impl) RemoveMember(ctx context.Context, orgID string, userID string) error {
	member, err := models.FindOrganizationMember(ctx, s.db, orgID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("find member: %w", err)
	}
	_, err = member.Delete(ctx, s.db)
	if err != nil {
		return fmt.Errorf("delete member: %w", err)
	}
	return nil
}

func (s *impl) UpdateMemberRole(ctx context.Context, orgID string, userID string, role string) error {
	m, err := models.FindOrganizationMember(ctx, s.db, orgID, userID)
	if err != nil {
		return fmt.Errorf("find member: %w", err)
	}
	m.Role = role
	_, err = m.Update(ctx, s.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("update member: %w", err)
	}
	return nil
}
