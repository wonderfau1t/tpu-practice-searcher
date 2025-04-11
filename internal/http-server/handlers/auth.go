package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage/models"
	"tpu-practice-searcher/internal/utils"
)

type AuthController interface {
	GetUserByID(userID int64) (*models.User, error)
}

type AuthResult struct {
	Role        string `json:"role"`
	AccessToken string `json:"accessToken"`
}

func Auth(log *slog.Logger, db AuthController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.Auth"
		log := log.With(slog.String("fn", fn))

		initData, ok := middlewares.CtxInitData(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Init data not found"))

			return
		}

		user, err := db.GetUserByID(initData.User.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, utils.NewErrorResponse("User not registered"))

				return
			}
			log.Error(fmt.Sprintf("failed to check existance of user: %s", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))

			return
		}

		accessToken, err := utils.GenerateAccessToken(user.Username, user.Role.Name)
		if err != nil {
			log.Error(fmt.Sprintf("failed to generate access token: %s", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))

			return
		}

		result := AuthResult{AccessToken: accessToken, Role: user.Role.Name}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(result))
	}
}
