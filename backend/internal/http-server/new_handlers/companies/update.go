package companies

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type UpdateCompanyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

type UpdateRepository interface {
	UpdateCompany(company *db_models.Company) error
	GetCompanyByHrID(hrID int64) (*db_models.HrManager, error)
	GetCompanyInfo(companyID uint) (*db_models.Company, error)
}

func Update(log *slog.Logger, db UpdateRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.Update"
		log := log.With(slog.String("fn", fn))

		_ = log

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			log.Info("not valid accessToken")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("failed to parse claims"))
			return
		}

		var req UpdateCompanyRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse body"))
			return
		}
		if req.Name == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("Field Name must not be empty"))
			return
		}

		companyID, err := db.GetCompanyByHrID(claims.UserID)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		company, err := db.GetCompanyInfo(companyID.CompanyID)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}
		updateCompany(company, req)
		err = db.UpdateCompany(company)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse("Company updated successfully"))
	}
}

func updateCompany(company *db_models.Company, req UpdateCompanyRequest) {
	company.Name = req.Name
	company.Description = utils.ToNullString(req.Description)
	company.Link = utils.ToNullString(req.Link)
}
