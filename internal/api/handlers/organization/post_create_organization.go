package organization

import (
	"net/http"

	"github.com/go-openapi/swag"
	"github.com/kashguard/go-mpc-vault/internal/api"
	"github.com/kashguard/go-mpc-vault/internal/auth"
	"github.com/kashguard/go-mpc-vault/internal/types"
	"github.com/kashguard/go-mpc-vault/internal/util"
	"github.com/labstack/echo/v4"
)

func PostCreateOrganizationRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Org.POST("", postCreateOrganizationHandler(s))
}

func postCreateOrganizationHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		u := auth.UserFromContext(ctx)
		var body types.CreateOrganizationPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}
		org, err := s.Organization.CreateOrganization(ctx, swag.StringValue(body.Name), u.ID)
		if err != nil {
			return err
		}
		return util.ValidateAndReturn(c, http.StatusOK, &types.CreateOrganizationResponse{
			OrganizationID: org.ID,
			Name:           org.Name,
		})
	}
}
