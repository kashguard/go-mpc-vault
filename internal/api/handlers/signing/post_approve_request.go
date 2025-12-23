package signing

import (
	"net/http"

	"github.com/go-openapi/swag"
	"github.com/kashguard/go-mpc-vault/internal/api"
	"github.com/kashguard/go-mpc-vault/internal/auth"
	"github.com/kashguard/go-mpc-vault/internal/service/signing"
	"github.com/kashguard/go-mpc-vault/internal/types"
	"github.com/kashguard/go-mpc-vault/internal/util"
	"github.com/labstack/echo/v4"
)

func PostApproveSigningRequestRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Sign.POST("/requests/:requestId/approve", postApproveSigningRequestHandler(s))
}

func postApproveSigningRequestHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		requestID := c.Param("requestId")

		var body types.ApproveSigningRequestPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		user := auth.UserFromContext(ctx)
		if user == nil {
			return echo.ErrUnauthorized
		}
		userID := user.ID

		action := swag.StringValue(body.Action)
		if action == "reject" {
			if err := s.Signing.RejectRequest(ctx, requestID, userID); err != nil {
				log.Error().Err(err).Msg("Failed to reject signing request")
				return err
			}
			return util.ValidateAndReturn(c, http.StatusOK, &types.ApproveSigningResponse{
				Status: "rejected",
			})
		}

		// Approve
		credID := []byte(body.CredentialID)
		sig := []byte(body.Signature)
		authData := []byte(body.AuthenticatorData)
		clientData := []byte(body.ClientDataJSON)

		params := signing.ApprovalParams{
			UserID:            userID,
			CredentialID:      credID,
			Signature:         sig,
			AuthenticatorData: authData,
			ClientDataJSON:    clientData,
		}

		if err := s.Signing.ApproveRequest(ctx, requestID, params); err != nil {
			log.Error().Err(err).Msg("Failed to approve signing request")
			return err
		}

		return util.ValidateAndReturn(c, http.StatusOK, &types.ApproveSigningResponse{
			Status: "approved",
		})
	}
}
