package organization

import (
	"net/http"

	"github.com/kashguard/go-mpc-vault/internal/api"
	"github.com/kashguard/go-mpc-vault/internal/types"
	"github.com/kashguard/go-mpc-vault/internal/util"
	"github.com/labstack/echo/v4"
)

func DeleteOrganizationMemberRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Org.DELETE("/:orgId/members/:userId", deleteOrganizationMemberHandler(s))
}

func deleteOrganizationMemberHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		orgID := c.Param("orgId")
		userID := c.Param("userId")
		if err := s.Organization.RemoveMember(ctx, orgID, userID); err != nil {
			return err
		}
		return util.ValidateAndReturn(c, http.StatusOK, &types.AddOrganizationMemberResponse{
			Ok: true,
		})
	}
}
