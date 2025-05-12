package vacancies

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/utils"
)

func Delete(log *slog.Logger, repo VacancyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers..Delete"
		log := log.With(slog.String("fn1", fn))

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			log.Info("not valid accessToken")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("failed to parse claims"))
			return
		}

		if claims.Role == "student" {
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, utils.NewErrorResponse(""))
			return
		}

		vacancyId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil || vacancyId == 0 {
			log.Info("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must be a valid positive integer"))
			return
		}
		err = repo.HideVacancyByID(uint(vacancyId))
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse("success"))
		return
	}
}
