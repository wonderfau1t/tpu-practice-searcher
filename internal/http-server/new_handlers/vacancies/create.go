package vacancies

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strings"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
	"tpu-practice-searcher/internal/utils/constants"
)

func Create(log *slog.Logger, repo VacancyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.Create"
		log := log.With(slog.String("fn", fn))

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse claims"))
			return
		}

		if claims.Role == "student" {
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, utils.NewErrorResponse("students cannot create vacancies"))
			return
		}

		user, err := repo.GetUserByID(claims.UserID)
		if err != nil {
			log.Error("failed to check phoneNumber")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("internal server error"))
			return
		}

		if !user.PhoneNumber.Valid {
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, utils.NewErrorResponse("phone number must be approved"))
			return
		}

		var req CreateVacancyRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse body"))
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				switch err.Tag() {
				case "required":
					errors = append(errors, err.Field()+" is required")
				default:
					errors = append(errors, err.Field()+" is invalid")
				}
			}
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("Validation failed: "+strings.Join(errors, ", ")))
			return
		}

		var vacancy *db_models.Vacancy
		switch claims.Role {
		case "HR", "headHR":
			vacancy = createVacancy(claims.UserID, claims.CompanyID, req)
		case "admin", "moderator":
			vacancy = createVacancyWithoutCompany(claims.UserID, req)
		}

		if err := repo.CreateVacancy(vacancy); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Failed to add vacancy"))
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, utils.NewSuccessResponse("success"))
	}
}

func createVacancy(userID int64, companyID uint, req CreateVacancyRequest) *db_models.Vacancy {
	vacancy := db_models.Vacancy{
		Name:                           req.VacancyName,
		CompanyID:                      &companyID,
		HrID:                           userID,
		StatusID:                       constants.StatusDefault,
		FormatID:                       req.FormatID,
		NumberOfResponses:              0,
		DeadlineAt:                     req.DeadlineAt,
		PaymentForAccommodationID:      req.PaymentForAccommodationID,
		PaymentForAccommodationDetails: utils.ToNullString(req.PaymentForAccommodationDetails),
		FarePaymentID:                  req.FarePaymentID,
		FarePaymentDetails:             utils.ToNullString(req.FarePaymentDetails),
		Description: db_models.VacancyDescription{
			Workplace:      utils.ToNullString(req.Description.Workplace),
			Position:       utils.ToNullString(req.Description.Position),
			Salary:         utils.ToNullString(req.Description.Salary),
			Requirements:   utils.ToNullString(req.Description.Requirements),
			Food:           utils.ToNullString(req.Description.Food),
			Conditions:     utils.ToNullString(req.Description.Conditions),
			AdditionalInfo: utils.ToNullString(req.Description.AdditionalInfo),
		},
	}

	for _, keyword := range req.Keywords {
		vacancy.Keywords = append(vacancy.Keywords, db_models.VacancyKeywords{Keyword: keyword})
	}
	for _, courseID := range req.Courses {
		vacancy.Courses = append(vacancy.Courses, db_models.Course{ID: courseID})
	}

	return &vacancy
}

func createVacancyWithoutCompany(userID int64, req CreateVacancyRequest) *db_models.Vacancy {
	vacancy := db_models.Vacancy{
		Name:                           req.VacancyName,
		CompanyName:                    &req.CompanyName,
		HrID:                           userID,
		StatusID:                       constants.StatusDefault,
		FormatID:                       req.FormatID,
		NumberOfResponses:              0,
		DeadlineAt:                     req.DeadlineAt,
		PaymentForAccommodationID:      req.PaymentForAccommodationID,
		PaymentForAccommodationDetails: utils.ToNullString(req.PaymentForAccommodationDetails),
		FarePaymentID:                  req.FarePaymentID,
		FarePaymentDetails:             utils.ToNullString(req.FarePaymentDetails),
		Description: db_models.VacancyDescription{
			Workplace:      utils.ToNullString(req.Description.Workplace),
			Position:       utils.ToNullString(req.Description.Position),
			Salary:         utils.ToNullString(req.Description.Salary),
			Requirements:   utils.ToNullString(req.Description.Requirements),
			Food:           utils.ToNullString(req.Description.Food),
			Conditions:     utils.ToNullString(req.Description.Conditions),
			AdditionalInfo: utils.ToNullString(req.Description.AdditionalInfo),
		},
	}

	for _, keyword := range req.Keywords {
		vacancy.Keywords = append(vacancy.Keywords, db_models.VacancyKeywords{Keyword: keyword})
	}
	for _, courseID := range req.Courses {
		vacancy.Courses = append(vacancy.Courses, db_models.Course{ID: courseID})
	}

	return &vacancy
}
