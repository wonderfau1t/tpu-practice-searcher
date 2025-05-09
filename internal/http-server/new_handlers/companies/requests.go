package companies

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type ReviewRepository interface {
	GetUnderReviewCompanies() ([]db_models.Company, error)
}

type CompanyDTO struct {
	CompanyID    uint   `json:"companyID"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Link         string `json:"link"`
	RegisteredAt string `json:"registeredAt"`
	HRUsername   string `json:"HRUsername"`
}

type RequestsResponse struct {
	TotalCount int          `json:"totalCount"`
	Companies  []CompanyDTO `json:"companies"`
}

func ListRequests(log *slog.Logger, db ReviewRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.Review"
		log := log.With(slog.String("fn", fn))

		_ = log

		companies, err := db.GetUnderReviewCompanies()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var companyDTOs []CompanyDTO
		for _, company := range companies {
			companyDTO := CompanyDTO{
				CompanyID:    company.ID,
				Name:         company.Name,
				Description:  company.Description.String,
				Link:         company.Link.String,
				RegisteredAt: company.CreatedAt.Format("2006-01-02"),
				HRUsername:   company.HeadHr.Username,
			}
			companyDTOs = append(companyDTOs, companyDTO)
		}

		response := RequestsResponse{
			TotalCount: len(companyDTOs),
			Companies:  companyDTOs,
		}

		render.JSON(w, r, utils.NewSuccessResponse(response))
		render.Status(r, http.StatusOK)
	}
}
