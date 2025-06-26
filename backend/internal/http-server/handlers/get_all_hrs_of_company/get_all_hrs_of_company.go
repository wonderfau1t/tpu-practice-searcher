package get_all_hrs_of_company

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/http-server/new_handlers/companies"
	"tpu-practice-searcher/internal/utils"
)

type HRDTO struct {
	Id               int64  `json:"id"`
	Username         string `json:"username"`
	CountOfVacancies int    `json:"countOfVacancies"`
}

type Response struct {
	TotalCount int               `json:"totalCount"`
	Hrs        []companies.HRDTO `json:"hrs"`
}

type Storage interface {
	GetAllHrsOfCompany(companyID uint) ([]companies.HRDTO, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.get_all_hrs_of_company.New"
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

		response := Response{
			TotalCount: len(hrs),
			Hrs:        hrs,
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(response))
	}
}
