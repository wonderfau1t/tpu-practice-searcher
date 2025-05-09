package companies

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"tpu-practice-searcher/internal/storage"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type DetailsRepository interface {
	GetCompanyInfo(companyID uint) (*db_models.Company, error)
}

func Details(log *slog.Logger, db DetailsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers..Details"
		log := log.With(slog.String("fn", fn))

		companyID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil || companyID == 0 {
			log.Info("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must be a valid positive integer"))
			return
		}

		company, err := db.GetCompanyInfo(uint(companyID))
		if err != nil {
			if errors.Is(err, storage.ErrRecordNotFound) {
				log.Info("vacancy not found")
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, utils.NewSuccessResponse("company not found"))
				return
			}
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("internal server error"))
			return
		}

		dto := CompanyDTO{
			CompanyID:    company.ID,
			Name:         company.Name,
			Description:  company.Description.String,
			Link:         company.Link.String,
			RegisteredAt: company.CreatedAt.Format("2006-01-02"),
			HRUsername:   company.HeadHr.Username,
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(dto))
	}
}
