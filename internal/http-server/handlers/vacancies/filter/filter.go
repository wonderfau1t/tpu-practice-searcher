package filter

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
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
	FilterVacancies(categoryID *uint, courseIDs []uint) ([]db_models.Vacancy, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.filter.New"
		log := log.With(slog.String("fn", fn))

		// Получаем category_id из query-параметров
		var categoryID *uint
		if catIDStr := r.URL.Query().Get("category_id"); catIDStr != "" {
			if catID, err := strconv.ParseUint(catIDStr, 10, 32); err == nil {
				catID32 := uint(catID)
				categoryID = &catID32
			} else {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, utils.NewErrorResponse("invalid category_id"))
				return
			}
		}

		// Получаем course_ids из query-параметров
		var courseIDs []uint
		if courseIDsStr := r.URL.Query().Get("course_ids"); courseIDsStr != "" {
			for _, idStr := range strings.Split(courseIDsStr, ",") {
				if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
					courseIDs = append(courseIDs, uint(id))
				} else {
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, utils.NewErrorResponse("invalid course_ids"))
					return
				}
			}
		}

		vacancies, err := db.FilterVacancies(categoryID, courseIDs)
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
				Category:    vacancy.Category.Name,
			}
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(Response{
			TotalCount: len(dtos),
			Vacancies:  dtos,
		}))
	}
}
