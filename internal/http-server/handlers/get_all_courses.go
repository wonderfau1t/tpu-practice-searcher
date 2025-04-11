package handlers

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/storage/models"
	"tpu-practice-searcher/internal/utils"
)

type GetAllCoursesController interface {
	GetAllCourses() ([]models.Course, error)
}

type GetAllCoursesResult struct {
	Courses []models.Course `json:"courses"`
}

func GetAllCourses(log *slog.Logger, db GetAllCoursesController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.GetAllCourses"
		log := log.With(slog.String("fn", fn))

		courses, err := db.GetAllCourses()
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		result := GetAllCoursesResult{Courses: courses}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(result))
	}
}
