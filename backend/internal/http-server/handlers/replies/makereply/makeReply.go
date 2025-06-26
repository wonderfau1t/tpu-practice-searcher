package makereply

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type Storage interface {
	GetUserByID(userID int64) (*db_models.User, error)
	AddReply(studentID int64, vacancyID uint) error
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse claims"))
			return
		}

		user, err := db.GetUserByID(claims.UserID)
		if err != nil {
			log.Error("failed to check phoneNumber")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("internal server error"))
			return
		}

		if !user.PhoneNumber.Valid {
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, utils.NewErrorResponse("phone number must be approved"))
			return
		}

		vacancyID, err := strconv.ParseUint(chi.URLParam(r, "vacancyID"), 10, 64)
		if err != nil || vacancyID == 0 {
			log.Info("Invalid vacancy ID")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must be a valid positive integer"))
			return
		}
		// FIX: Запрос может прилететь на несуществующую вакансию

		if err := db.AddReply(claims.UserID, uint(vacancyID)); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, utils.NewSuccessResponse("replied!"))
	}
}
