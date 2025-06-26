package companies

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/utils"
)

type Request struct {
	Username string `json:"username"`
}

func CreateHr(log *slog.Logger, db CompanyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.CreateHR"
		log := log.With(slog.String("fn", fn))

		_ = log

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse claims"))
			return
		}

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse body"))
			return
		}

		if req.Username == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("Name must not be empty"))
			return
		}

		if err := db.CreateNewHr(req.Username, claims.CompanyID); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("InternalServerError"))
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse("Successfully registered"))
	}
}
