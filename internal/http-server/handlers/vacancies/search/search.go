package search

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strings"
	common "tpu-practice-searcher/internal/http-server/handlers/vacancies"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type Response struct {
	TotalCount int                    `json:"totalCount"`
	Vacancies  []common.GetVacancyDTO `json:"vacancies"`
}

type Storage interface {
	SearchVacancies(searchQuery string) ([]db_models.Vacancy, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.search.New"
		log := log.With(slog.String("fn", fn))

		searchQuery := strings.TrimSpace(r.URL.Query().Get("query"))
		if searchQuery == "" {
			log.Info("query is not set")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("query is not set"))
			return
		}

		vacancies, err := db.SearchVacancies(searchQuery)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		dtos := make([]common.GetVacancyDTO, len(vacancies))
		for i, vacancy := range vacancies {
			dtos[i] = common.GetVacancyDTO{
				ID:          vacancy.ID,
				Name:        vacancy.Name,
				CompanyName: vacancy.Company.Name,
			}
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(Response{
			TotalCount: len(dtos),
			Vacancies:  dtos,
		}))
	}
}
