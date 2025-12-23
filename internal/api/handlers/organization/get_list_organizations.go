package organization

import (
	"net/http"

	"github.com/kashguard/go-mpc-vault/internal/api"
	"github.com/kashguard/go-mpc-vault/internal/auth"
	"github.com/kashguard/go-mpc-vault/internal/types"
	"github.com/kashguard/go-mpc-vault/internal/util"
	"github.com/labstack/echo/v4"
)

func GetListOrganizationsRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Org.GET("", getListOrganizationsHandler(s))
}

func getListOrganizationsHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		u := auth.UserFromContext(ctx)
		orgs, err := s.Organization.ListUserOrganizations(ctx, u.ID)
		if err != nil {
			return err
		}
		items := make([]*types.OrganizationItem, 0, len(orgs))
		for _, o := range orgs {
			items = append(items, &types.OrganizationItem{
				ID:   o.ID,
				Name: o.Name,
			})
		}
		return util.ValidateAndReturn(c, http.StatusOK, &types.ListOrganizationsResponse{
			Organizations: items,
		})
	}
}
