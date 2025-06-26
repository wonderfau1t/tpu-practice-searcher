package updatecompanyinfo

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/utils"
)

type Storage interface {
	UpdateCompanyInfo(companyID uint, name string, description string, link string) error
}

type Request struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.updatecompanyinfo.New"
		log := log.With(slog.String("fn", fn))

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse claims"))
			return
		}

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("failed to parse body"))
			return
		}

		if req.Name == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("field name must not be empty"))
			return
		}

		err := db.UpdateCompanyInfo(claims.CompanyID, req.Name, req.Description, req.Link)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse("success"))
	}
}
