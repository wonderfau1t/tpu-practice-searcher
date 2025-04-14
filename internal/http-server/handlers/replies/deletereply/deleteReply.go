package deletereply

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/utils"
)

type Storage interface {
	DeleteReply(studentID int64, vacancyID uint) error
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.deletereply.New"
		log := log.With(slog.String("fn", fn))

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse claims"))
			return
		}

		vacancyID, err := strconv.ParseUint(chi.URLParam(r, "vacancyID"), 10, 64)
		if err != nil || vacancyID == 0 {
			log.Info("Invalid vacancy ID")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must be a valid positive integer"))
			return
		}

		if err := db.DeleteReply(claims.UserID, uint(vacancyID)); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, utils.NewSuccessResponse("reply deleted"))
	}
}
