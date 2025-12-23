package signing

import (
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/kashguard/go-mpc-vault/internal/api"
	"github.com/kashguard/go-mpc-vault/internal/auth"
	"github.com/kashguard/go-mpc-vault/internal/types"
	"github.com/kashguard/go-mpc-vault/internal/util"
	"github.com/labstack/echo/v4"
)

func PostCreateSigningRequestRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Sign.POST("/vaults/:vaultId/sign", postCreateSigningRequestHandler(s))
}

func postCreateSigningRequestHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		// vaultID := c.Param("vaultId") // Unused in service create request, but good for validation

		var body types.CreateSigningRequestPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		user := auth.UserFromContext(ctx)
		if user == nil {
			return echo.ErrUnauthorized
		}
		userID := user.ID

		req, err := s.Signing.CreateRequest(ctx, string(*body.WalletID), body.TxData, body.Note, userID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create signing request")
			return err
		}

		return util.ValidateAndReturn(c, http.StatusOK, &types.CreateSigningResponse{
			RequestID: strfmt.UUID4(req.ID),
			Status:    req.Status.String,
		})
	}
}
