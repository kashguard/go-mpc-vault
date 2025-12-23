package organization

import (
	"net/http"

	"github.com/go-openapi/swag"
	"github.com/kashguard/go-mpc-vault/internal/api"
	"github.com/kashguard/go-mpc-vault/internal/types"
	"github.com/kashguard/go-mpc-vault/internal/util"
	"github.com/labstack/echo/v4"
)

func PostAddOrganizationMemberRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Org.POST("/:orgId/members", postAddOrganizationMemberHandler(s))
}

func postAddOrganizationMemberHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		orgID := c.Param("orgId")
		var body types.AddOrganizationMemberPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}
		_, err := s.Organization.AddMember(ctx, orgID, swag.StringValue(body.UserID), swag.StringValue(body.Role))
		if err != nil {
			return err
		}
		return util.ValidateAndReturn(c, http.StatusOK, &types.AddOrganizationMemberResponse{
			Ok: true,
		})
	}
}
