package get_all_vacancies_by_courses

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type VacancyDTO struct {
	Id             uint   `json:"id"`
	Name           string `json:"name"`
	CompanyName    string `json:"companyName"`
	Category       string `json:"category"`
	CountOfReplies int    `json:"countOfReplies"`
}

type Response struct {
	TotalCount int          `json:"totalCount"`
	Vacancies  []VacancyDTO `json:"vacancies"`
}

type Storage interface {
	GetSchoolByModeratorID(moderatorID int64) (uint, error)
	GetVacanciesBySchoolID(schoolID uint) ([]db_models.Vacancy, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse claims"))
			return
		}
		schoolID, err := db.GetSchoolByModeratorID(claims.UserID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, utils.NewErrorResponse("Failed to find School"))
			return
		}

		vacancies, err := db.GetVacanciesBySchoolID(schoolID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, utils.NewErrorResponse("Failed to find vacancies"))
			return
		}

		dtos := make([]VacancyDTO, len(vacancies))
		for i, vacancy := range vacancies {
			dtos[i] = VacancyDTO{
				Id:             vacancy.ID,
				Name:           vacancy.Name,
				CompanyName:    vacancy.Company.Name,
				CountOfReplies: vacancy.NumberOfResponses,
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
