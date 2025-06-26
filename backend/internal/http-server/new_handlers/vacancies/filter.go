package vacancies

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	common "tpu-practice-searcher/internal/http-server/handlers/vacancies"
	"tpu-practice-searcher/internal/utils"
)

func Filter(log *slog.Logger, repo VacancyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.Search"
		log := log.With(slog.String("fn", fn))

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

		vacancies, err := repo.FilterVacancies(courseIDs)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		dtos := make([]common.GetVacancyDTO, len(vacancies))
		for i, vacancy := range vacancies {
			var companyName string
			if vacancy.CompanyID == nil {
				companyName = *vacancy.CompanyName
			} else {
				companyName = vacancy.Company.Name
			}
			dtos[i] = common.GetVacancyDTO{
				ID:          vacancy.ID,
				Name:        vacancy.Name,
				CompanyName: companyName,
			}
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(Response{
			TotalCount: len(dtos),
			Vacancies:  dtos,
		}))
	}
}
