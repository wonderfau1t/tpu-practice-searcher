package handlers

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type GetAllFormatsController interface {
	GetAllFormats() ([]db_models.Format, error)
}

type GetAllFormatsResult struct {
	Formats []db_models.Format `json:"formats"`
}

func GetAllFormats(log *slog.Logger, db GetAllFormatsController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.GetAllFormats"
		log := log.With(slog.String("fn", fn))

		formats, err := db.GetAllFormats()
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		result := GetAllFormatsResult{Formats: formats}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(result))
	}
}
