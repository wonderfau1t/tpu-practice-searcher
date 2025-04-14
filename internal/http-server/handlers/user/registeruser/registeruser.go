package registeruser

import (
	"fmt"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/utils"
	"tpu-practice-searcher/internal/utils/constants"
)

type Storage interface {
	IsUserExist(userId int64) (bool, error)
	CreateNewUser(userId int64, username string, roleId uint) error
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.registeruser.New"
		log := log.With(slog.String("fn", fn))

		initData, ok := middlewares.CtxInitData(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("initData is not valid"))
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
		err = db.CreateNewUser(initData.User.ID, initData.User.Username, constants.RoleStudent)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to register new unser: %s", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		// FIX: Уйти от хардкодинга роли
		accessToken, err := utils.GenerateStudentAccessToken(initData.User.ID, initData.User.Username, "student")
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
