package get_company_details

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type CompanyDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
}

type Response struct {
	CompanyInfo CompanyDTO `json:"companyInfo"`
}

type Storage interface {
	GetCompanyInfo(companyID uint) (*db_models.Company, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse claims"))
			return
		}

		company, err := db.GetCompanyInfo(claims.CompanyID)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		dto := CompanyDTO{
			Name:        company.Name,
			Description: company.Description.String,
			Website:     company.Link.String,
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(dto))
	}
}
