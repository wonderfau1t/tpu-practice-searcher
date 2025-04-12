package get_all_vacancies_of_company

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type VacancyDTO struct {
	ID               uint     `json:"id"`
	Name             string   `json:"name"`
	Courses          []string `json:"courses"`
	Category         string   `json:"category"`
	CountOfResponses int      `json:"responses"`
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

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse claims"))
			return
		}

		exists, err := db.IsCompanyExist(claims.CompanyID)
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

		vacancies, err := db.GetAllVacanciesOfCompany(claims.CompanyID)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse(err.Error()))
			return
		}

		dtos := make([]VacancyDTO, len(vacancies))
		for i, vacancy := range vacancies {
			dtos[i] = VacancyDTO{
				ID:               vacancy.ID,
				Name:             vacancy.Name,
				Category:         vacancy.Category.Name,
				CountOfResponses: vacancy.NumberOfResponses,
			}
			for _, course := range vacancy.Courses {
				dtos[i].Courses = append(dtos[i].Courses, course.Name)
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
