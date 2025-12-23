package signing

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/kashguard/go-mpc-vault/internal/api"
	"github.com/kashguard/go-mpc-vault/internal/auth"
	"github.com/kashguard/go-mpc-vault/internal/models"
	"github.com/kashguard/go-mpc-vault/internal/util"
	"github.com/labstack/echo/v4"
)

type ListSigningRequestsParams struct {
	Status  string `query:"status"`
	VaultID string `query:"vaultId"`
	Page    int    `query:"page"`
	Limit   int    `query:"limit"`
}

func (p *ListSigningRequestsParams) Validate(_ strfmt.Registry) error {
	switch p.Status {
	case "", "pending", "completed", "rejected":
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid status")
	}
	if p.Page < 0 || p.Limit < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid pagination")
	}
	return nil
}

type SigningRequestItem struct {
	ID        string `json:"id"`
	VaultID   string `json:"vault_id,omitempty"`
	WalletID  string `json:"wallet_id,omitempty"`
	ToAddress string `json:"to_address,omitempty"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at,omitempty"`
}

func (_ *SigningRequestItem) Validate(_ strfmt.Registry) error { return nil }

type ListSigningRequestsResponse struct {
	Requests []*SigningRequestItem `json:"requests"`
	Total    int32                 `json:"total"`
}

func (_ *ListSigningRequestsResponse) Validate(_ strfmt.Registry) error { return nil }

func GetListSigningRequestsRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1Sign.GET("/requests", getListSigningRequestsHandler(s))
}

func getListSigningRequestsHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		var params ListSigningRequestsParams
		if err := util.BindAndValidateQueryParams(c, runtime.Validatable(&params)); err != nil {
			return err
		}

		user := auth.UserFromContext(ctx)
		if user == nil {
			return echo.ErrUnauthorized
		}
		userID := user.ID

		items, total, err := s.Signing.ListRequests(ctx, userID, params.VaultID, params.Status, params.Page, params.Limit)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list signing requests")
			return err
		}

		resp := &ListSigningRequestsResponse{
			Requests: make([]*SigningRequestItem, 0, len(items)),
			Total:    int32(total),
		}

		for _, r := range items {
			resp.Requests = append(resp.Requests, mapSigningRequest(r))
		}

		return util.ValidateAndReturn(c, http.StatusOK, resp)
	}
}

func mapSigningRequest(r *models.SigningRequest) *SigningRequestItem {
	item := &SigningRequestItem{
		ID:     r.ID,
		Status: r.Status.String,
	}
	if r.VaultID.Valid {
		item.VaultID = r.VaultID.String
	}
	if r.WalletID.Valid {
		item.WalletID = r.WalletID.String
	}
	if r.ToAddress.Valid {
		item.ToAddress = r.ToAddress.String
	}
	if r.CreatedAt.Valid {
		item.CreatedAt = r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}
	return item
}
