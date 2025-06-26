package companies

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/utils"
)

type response struct {
	TotalCount int     `json:"totalCount"`
	Hrs        []HRDTO `json:"hrs"`
}

func HrList(log *slog.Logger, db CompanyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.HrList"
		log := log.With(slog.String("fn", fn))

		_ = log

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Unauthorized"))
			return
		}

		hrs, err := db.GetAllHrsOfCompany(claims.CompanyID)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		response := response{
			TotalCount: len(hrs),
			Hrs:        hrs,
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(response))
	}
}
