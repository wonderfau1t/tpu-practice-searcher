package createvacancywithoutcompany

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

type VacancyDescription struct {
	Workplace      string `json:"workplace"`
	Position       string `json:"position"`
	Salary         string `json:"salary"`
	Requirements   string `json:"requirements"`
	Food           string `json:"food"`
	Conditions     string `json:"conditions"`
	AdditionalInfo string `json:"additionalInfo"`
}

type Request struct {
	CompanyName                    string             `json:"companyName" validate:"required"`
	VacancyName                    string             `json:"vacancyName" validate:"required"`
	FormatID                       uint               `json:"formatID" validate:"required"`
	Courses                        []uint             `json:"courses" validate:"required"`
	Keywords                       []string           `json:"keywords"`
	DeadlineAt                     string             `json:"deadlineAt" validate:"required"`
	PaymentForAccommodationID      uint               `json:"paymentForAccommodationID" validate:"required"`
	PaymentForAccommodationDetails string             `json:"paymentForAccommodationDetails"`
	FarePaymentID                  uint               `json:"farePaymentID" validate:"required"`
	FarePaymentDetails             string             `json:"farePaymentDetails"`
	Description                    VacancyDescription `json:"description,omitempty"`
}

type Storage interface {
	GetUserByID(userID int64) (*db_models.User, error)
	CreateNewVacancy(vacancy *db_models.Vacancy) error
	GetCompanyByHrID(hrID int64) (*db_models.HrManager, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.createvacancywithoutcompany.New"
		log := log.With(slog.String("fn", fn))

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse claims"))
			return
		}

		user, err := db.GetUserByID(claims.UserID)
		if err != nil {
			log.Error("failed to check phoneNumber")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("internal server error"))
			return
		}
		_ = user

		// Должен ли быть подтвержденным номер телефона у админа или модератора
		//if !user.PhoneNumber.Valid {
		//	render.Status(r, http.StatusForbidden)
		//	render.JSON(w, r, utils.NewErrorResponse("phone number must be approved"))
		//	return
		//}

		var req Request
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

		vacancy := db_models.Vacancy{
			Name:                           req.VacancyName,
			CompanyName:                    &req.CompanyName,
			HrID:                           claims.UserID,
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

		if err := db.CreateNewVacancy(&vacancy); err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Failed to add vacancy"))
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, utils.NewSuccessResponse("success"))
	}
}
