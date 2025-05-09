package companies

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/utils"
	"tpu-practice-searcher/internal/utils/constants"
)

type RejectCompanyRequest struct {
	CompanyID uint   `json:"companyID"`
	Message   string `json:"message"`
}

type RejectRepository interface {
	UpdateCompanyStatus(companyID uint, statusID int) error
}

func Reject(log *slog.Logger, repo RejectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.Apply"
		log := log.With(slog.String("fn", fn))

		var req ApplyCompanyRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse body"))
			return
		}

		err := repo.UpdateCompanyStatus(req.CompanyID, constants.StatusBlocked)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse("Company status updated successfully"))
	}
}
