package get_vacancies

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
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
	GetAllVacancies() ([]db_models.Vacancy, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vacancies, err := db.GetAllVacancies()
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		dtos := make([]VacancyDTO, len(vacancies))

		for i, vacancy := range vacancies {
			dtos[i] = VacancyDTO{
				Id:          vacancy.ID,
				Name:        vacancy.Name,
				CompanyName: vacancy.Company.Name,
				Category:    vacancy.Category.Name,
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
