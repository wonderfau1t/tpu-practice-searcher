package handlers

import (
	"database/sql"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/utils"
)

type RegisterCompanyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

type RegisterCompanyController interface {
	CreateNewCompany(userId int64, username string, companyName string, companyDescription sql.NullString, companyLink sql.NullString) error
}

func RegisterCompany(log *slog.Logger, db RegisterCompanyController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.RegisterCompany"
		log := log.With(slog.String("fn", fn))

		initData, ok := middlewares.CtxInitData(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse initData"))
			return
		}

		//initData.User.ID = 1337
		//initData.User.Username = "wonderrfau1t"

		var req RegisterCompanyRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse body"))
			return
		}

		if req.Name == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("Field Name must not be empty"))
		}

		description := utils.ToNullString(req.Description)
		link := utils.ToNullString(req.Link)

		err = db.CreateNewCompany(initData.User.ID, initData.User.Username, req.Name, description, link)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, utils.NewSuccessResponse("CompanyCreated"))
	}
}
