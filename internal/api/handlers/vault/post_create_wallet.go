package vault

import (
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/kashguard/go-mpc-vault/internal/api"
	"github.com/kashguard/go-mpc-vault/internal/auth"
	"github.com/kashguard/go-mpc-vault/internal/models"
	"github.com/kashguard/go-mpc-vault/internal/types"
	"github.com/kashguard/go-mpc-vault/internal/util"
	"github.com/labstack/echo/v4"
)

func PostCreateWalletRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Vault.POST("/:vaultId/wallets", postCreateWalletHandler(s))
}

func postCreateWalletHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		vaultID := c.Param("vaultId")

		var body types.CreateWalletPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		user := auth.UserFromContext(ctx)
		if user == nil {
			return echo.ErrUnauthorized
		}
		v, err := models.FindVault(ctx, s.DB, vaultID)
		if err != nil {
			return echo.ErrNotFound
		}
		ok, err := models.OrganizationMembers(
			models.OrganizationMemberWhere.OrganizationID.EQ(v.OrganizationID.String),
			models.OrganizationMemberWhere.UserID.EQ(user.ID),
		).Exists(ctx, s.DB)
		if err != nil {
			return err
		}
		if !ok && v.OrganizationID.Valid {
			return echo.ErrForbidden
		}

		wallet, err := s.Vault.CreateWallet(ctx, vaultID, swag.StringValue(body.ChainID))
		if err != nil {
			log.Error().Err(err).Msg("Failed to create wallet")
			return err
		}

		return util.ValidateAndReturn(c, http.StatusOK, &types.CreateWalletResponse{
			WalletID: strfmt.UUID4(wallet.ID),
			Address:  wallet.Address,
			ChainID:  wallet.ChainID.String,
		})
	}
}
