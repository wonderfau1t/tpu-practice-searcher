package handlers

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type GetAllCategoriesController interface {
	GetAllCategories() ([]db_models.Category, error)
}

type GetAllCategoriesResult struct {
	Categories []db_models.Category `json:"categories"`
}

func GetAllCategories(log *slog.Logger, db GetAllCategoriesController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.GetAllCategories"
		log := log.With(slog.String("fn", fn))

		categories, err := db.GetAllCategories()
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		result := GetAllCategoriesResult{Categories: categories}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(result))
	}
}
