package get_company_info

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
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
	IsCompanyExist(companyID uint) (bool, error)
	GetCompanyInfo(companyID uint) (*db_models.Company, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		companyIdStr := chi.URLParam(r, "id")
		if companyIdStr == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must not be empty"))
			return
		}
		companyId, err := strconv.ParseUint(companyIdStr, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must be type of uint"))
			return
		}

		exists, err := db.IsCompanyExist(uint(companyId))
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, utils.NewErrorResponse("company not found"))
			return
		}

		if !exists {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, utils.NewErrorResponse("company not found"))
			return
		}

		company, err := db.GetCompanyInfo(uint(companyId))
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
