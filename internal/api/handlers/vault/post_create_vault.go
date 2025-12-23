package vault

import (
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/kashguard/go-mpc-vault/internal/api"
	"github.com/kashguard/go-mpc-vault/internal/auth"
	"github.com/kashguard/go-mpc-vault/internal/types"
	"github.com/kashguard/go-mpc-vault/internal/util"
	"github.com/labstack/echo/v4"
)

func PostCreateVaultRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Vault.POST("", postCreateVaultHandler(s))
}

func postCreateVaultHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var body types.CreateVaultPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		user := auth.UserFromContext(ctx)
		if user == nil {
			return echo.ErrUnauthorized
		}
		orgs, err := s.Organization.ListUserOrganizations(ctx, user.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list user organizations")
			return err
		}
		if len(orgs) == 0 {
			return echo.ErrForbidden
		}
		orgID := orgs[0].ID
		for _, o := range orgs {
			if o.OwnerID == user.ID {
				orgID = o.ID
				break
			}
		}

		vault, err := s.Vault.CreateVault(ctx, swag.StringValue(body.Name), orgID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create vault")
			return err
		}

		if len(body.Chains) > 0 {
			for _, chainID := range body.Chains {
				_, err := s.Vault.CreateWallet(ctx, vault.ID, chainID)
				if err != nil {
					log.Error().Err(err).Str("chain_id", chainID).Msg("Failed to create initial wallet")
					// Continue or fail?
				}
			}
		}

		return util.ValidateAndReturn(c, http.StatusOK, &types.CreateVaultResponse{
			VaultID: strfmt.UUID4(vault.ID),
			Status:  "active",
		})
	}
}
