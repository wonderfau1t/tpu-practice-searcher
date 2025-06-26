package get_vacancies_of_company

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type VacancyDTO struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	CompanyName string `json:"companyName"`
	Category    string `json:"category"`
}

type Response struct {
	TotalCount int          `json:"totalCount"`
	Vacancies  []VacancyDTO `json:"vacancies"`
}

type Storage interface {
	IsCompanyExist(companyID uint) (bool, error)
	GetAllVacanciesOfCompany(companyID uint) ([]db_models.Vacancy, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.get_all_vacancies.New"
		_ = log.With(slog.String("fn", fn))

		companyIdStr := chi.URLParam(r, "id")
		if companyIdStr == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must not be empty"))
			return
		}
		companyId, err := strconv.ParseUint(companyIdStr, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must be type of uint"))
			return
		}

		exists, err := db.IsCompanyExist(uint(companyId))
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, utils.NewErrorResponse("company not found"))
			return
		}

		if !exists {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, utils.NewErrorResponse("company not found"))
			return
		}

		vacancies, err := db.GetAllVacanciesOfCompany(uint(companyId))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse(err.Error()))
			return
		}

		dtos := make([]VacancyDTO, len(vacancies))
		for i, vacancy := range vacancies {
			dtos[i] = VacancyDTO{
				Id:          vacancy.ID,
				Name:        vacancy.Name,
				CompanyName: vacancy.Company.Name,
			}
		}

		response := Response{
			TotalCount: len(dtos),
			Vacancies:  dtos,
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(response))
	}
}
