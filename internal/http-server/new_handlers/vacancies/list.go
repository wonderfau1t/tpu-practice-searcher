package vacancies

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type VacancyDTO struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	CompanyName    string `json:"companyName,omitempty"`
	CountOfReplies *int   `json:"countOfReplies,omitempty"`
}

type Response struct {
	TotalCount int          `json:"totalCount"`
	Vacancies  []VacancyDTO `json:"vacancies"`
}

type VacancyRepository interface {
	GetAllVacancies() ([]db_models.Vacancy, error)
}

func List(log *slog.Logger, repo VacancyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.List"
		log := log.With(slog.String("fn", fn))
		_ = log

		vacancies, err := repo.GetAllVacancies()
		if err != nil {
			return
		}

		dtos := make([]VacancyDTO, len(vacancies))
		for i, vacancy := range vacancies {

			var companyName string
			if vacancy.CompanyID == nil {
				companyName = *vacancy.CompanyName
			} else {
				companyName = vacancy.Company.Name
			}

			dtos[i] = VacancyDTO{
				ID:          vacancy.ID,
				Name:        vacancy.Name,
				CompanyName: companyName,
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
