package companies

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"tpu-practice-searcher/internal/utils"
	"tpu-practice-searcher/internal/utils/constants"
)

type RejectCompanyRequest struct {
	Message string `json:"message"`
}

type RejectRepository interface {
	UpdateCompanyStatus(companyID uint, statusID int) error
}

func Reject(log *slog.Logger, repo RejectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.Apply"
		log := log.With(slog.String("fn", fn))

		companyID, err := strconv.ParseUint(chi.URLParam(r, "companyID"), 10, 64)
		if err != nil || companyID == 0 {
			log.Info("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must be a valid positive integer"))
			return
		}

		err = repo.UpdateCompanyStatus(uint(companyID), constants.StatusDefault)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		err = repo.UpdateCompanyStatus(uint(companyID), constants.StatusBlocked)
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
