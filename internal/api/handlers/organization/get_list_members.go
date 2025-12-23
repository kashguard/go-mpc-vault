package organization

import (
	"net/http"

	"github.com/kashguard/go-mpc-vault/internal/api"
	"github.com/kashguard/go-mpc-vault/internal/types"
	"github.com/kashguard/go-mpc-vault/internal/util"
	"github.com/labstack/echo/v4"
)

func GetListOrganizationMembersRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Org.GET("/:orgId/members", getListOrganizationMembersHandler(s))
}

func getListOrganizationMembersHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		orgID := c.Param("orgId")
		members, err := s.Organization.ListMembers(ctx, orgID)
		if err != nil {
			return err
		}
		items := make([]*types.OrganizationMemberItem, 0, len(members))
		for _, m := range members {
			items = append(items, &types.OrganizationMemberItem{
				UserID: m.UserID,
				Role:   m.Role,
			})
		}
		return util.ValidateAndReturn(c, http.StatusOK, &types.ListOrganizationMembersResponse{
			Members: items,
		})
	}
}
