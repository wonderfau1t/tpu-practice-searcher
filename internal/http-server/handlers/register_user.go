package handlers

import (
	"fmt"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/utils"
)

type RegisterStudentController interface {
	IsUserExist(userId int64) (bool, error)
	CreateNewUser(userId int64, username string, roleId uint) error
}

type RegisterResult struct {
	Role        string `json:"role"`
	AccessToken string `json:"accessToken"`
}

func RegisterStudent(log *slog.Logger, db RegisterStudentController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.RegisterStudent"

		log := log.With(slog.String("fn", fn))

		initData, ok := middlewares.CtxInitData(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Init data not found"))

			return
		}

		exists, err := db.IsUserExist(initData.User.ID)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to check existence of user: %s", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))

			return
		}

		if exists {
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, utils.NewErrorResponse("User already exists"))

			return
		}

		// FIX: Уйти от хардкодинга роли
		err = db.CreateNewUser(initData.User.ID, initData.User.Username, 4)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to register new unser: %s", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))

			return
		}

		// FIX: Уйти от хардкодинга роли
		accessToken, err := utils.GenerateAccessToken(initData.User.Username, "student")
		if err != nil {
			log.Error(fmt.Sprintf("failed to generate access token: %s", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))

			return
		}

		result := RegisterResult{AccessToken: accessToken, Role: "student"}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(result))
	}
}
