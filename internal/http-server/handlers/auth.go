package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type AuthController interface {
	//GetUserByID(userID int64) (*db_models.User, error)
	GetCompanyIDByHRID(hrID int64) (uint, error)
	//IsRegisteredByUsername(username string) (bool, error)
	//RegisterByUsername(hrID int64) error
	FindOrRegisterByUsername(userID int64, username string) (*db_models.User, bool, error)
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

		user, registered, err := db.FindOrRegisterByUsername(initData.User.ID, initData.User.Username)
		if err != nil {
			if errors.Is(err, storage.ErrRecordNotFound) {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, utils.NewErrorResponse("not registered"))
				return
			}
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("internal server error"))
			return
		}

		if !registered {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("not registered"))
			return
		}

		var accessToken string
		switch user.Role.Name {
		case "student":
			accessToken, err = utils.GenerateStudentAccessToken(user.ID, user.Username, user.Role.Name)
		case "moderator":
			accessToken, err = utils.GenerateStudentAccessToken(user.ID, user.Username, user.Role.Name)
		case "HR", "headHR":
			companyID, err := db.GetCompanyIDByHRID(user.ID)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
				return
			}
			accessToken, err = utils.GenerateHrAccessToken(user.ID, user.Username, companyID, user.Role.Name)
		}
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
