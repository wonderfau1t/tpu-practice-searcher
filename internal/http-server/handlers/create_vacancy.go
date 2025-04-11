package handlers

import (
	"database/sql"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"tpu-practice-searcher/internal/storage/models"
	"tpu-practice-searcher/internal/utils"
	"tpu-practice-searcher/internal/utils/constants"
)

type VacancyDescription struct {
	Workplace      string `json:"workplace"`
	Position       string `json:"position"`
	Salary         int64  `json:"salary"`
	Requirements   string `json:"requirements"`
	Food           string `json:"food"`
	Conditions     string `json:"conditions"`
	AdditionalInfo string `json:"additionalInfo"`
}

type AddVacancyRequest struct {
	Name                           string             `json:"name" validate:"required"`
	FormatID                       uint               `json:"formatID" validate:"required"`
	CategoryID                     uint               `json:"categoryID" validate:"required"`
	Courses                        []uint             `json:"courses" validate:"required"`
	Keywords                       []string           `json:"keywords"`
	DeadlineAt                     time.Time          `json:"deadlineAt" validate:"required"`
	PaymentForAccommodationID      uint               `json:"paymentForAccommodationID" validate:"required"`
	PaymentForAccommodationDetails string             `json:"paymentForAccommodationDetails"`
	FarePaymentID                  uint               `json:"farePaymentID" validate:"required"`
	FarePaymentDetails             string             `json:"farePaymentDetails"`
	Description                    VacancyDescription `json:"description"`
}

type AddVacancyResult struct {
	ID uint `json:"ID"`
}

type AddVacancyController interface {
	CreateNewVacancy(vacancy *models.Vacancy) error
}

func AddVacancy(log *slog.Logger, db AddVacancyController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.AddVacancy"
		log := log.With(slog.String("fn", fn))

		var req AddVacancyRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse body"))
			return
		}

		validate := validator.New()
		//if err := validate.Struct(req); err != nil {
		//	render.Status(r, http.StatusBadRequest)
		//	render.JSON(w, r, utils.NewErrorResponse("Required fields must not be empty"))
		//	return
		//}

		if err := validate.Struct(req); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				// Формируем читаемое сообщение об ошибке
				switch err.Tag() {
				case "required":
					errors = append(errors, err.Field()+" is required")
				case "gt":
					errors = append(errors, err.Field()+" must be greater than "+err.Param())
				case "oneof":
					errors = append(errors, err.Field()+" must be one of "+err.Param())
				default:
					errors = append(errors, err.Field()+" is invalid")
				}
			}
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("Validation failed: "+strings.Join(errors, ", ")))
			return
		}

		vacancy := models.Vacancy{
			Name: req.Name,
			// Откуда брать ID компании
			CompanyID:                      0,
			HrID:                           0,
			StatusID:                       constants.StatusDefault,
			FormatID:                       req.FormatID,
			CategoryID:                     req.CategoryID,
			NumberOfResponses:              0,
			DeadlineAt:                     req.DeadlineAt,
			PaymentForAccommodationID:      req.PaymentForAccommodationID,
			PaymentForAccommodationDetails: utils.ToNullString(req.PaymentForAccommodationDetails),
			FarePaymentID:                  req.FarePaymentID,
			FarePaymentDetails:             utils.ToNullString(req.FarePaymentDetails),
			Description: models.VacancyDescription{
				Workplace: utils.ToNullString(req.Description.Workplace),
				Position:  utils.ToNullString(req.Description.Position),
				Salary: sql.NullInt64{
					Int64: req.Description.Salary,
					Valid: true,
				},
				Requirements:   utils.ToNullString(req.Description.Requirements),
				Food:           utils.ToNullString(req.Description.Food),
				Conditions:     utils.ToNullString(req.Description.Conditions),
				AdditionalInfo: utils.ToNullString(req.Description.AdditionalInfo),
			},
		}

		for _, keyword := range req.Keywords {
			vacancy.Keywords = append(vacancy.Keywords, models.VacancyKeywords{Keyword: keyword})
		}
		for _, courseID := range req.Courses {
			vacancy.Courses = append(vacancy.Courses, models.Course{ID: courseID})
		}

		if err := db.CreateNewVacancy(&vacancy); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Failed to add vacancy"))
			return
		}

		result := AddVacancyResult{ID: vacancy.ID}
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, utils.NewSuccessResponse(result))
	}
}
