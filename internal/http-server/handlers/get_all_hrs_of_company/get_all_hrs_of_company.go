package get_all_hrs_of_company

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type HRDTO struct {
	Id               int64  `json:"id"`
	Username         string `json:"username"`
	CountOfVacancies string `json:"countOfVacancies"`
}

type Response struct {
	TotalCount int     `json:"totalCount"`
	Hrs        []HRDTO `json:"hrs"`
}

type Storage interface {
	GetAllHrsOfCompany(companyID uint) ([]db_models.User, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.get_all_hrs_of_company.New"
		log := log.With(slog.String("fn", fn))

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Unauthorized"))
			return
		}

	}
}
