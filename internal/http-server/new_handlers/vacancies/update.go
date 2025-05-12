package vacancies

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type CoursesRequest struct {
	CourseID uint   `json:"courseId"`
	Name     string `json:"name"`
}

type UpdateVacancyRequest struct {
	CompanyName                    string             `json:"companyName,omitempty"`
	VacancyName                    string             `json:"vacancyName" validate:"required"`
	FormatID                       uint               `json:"formatID" validate:"required"`
	Courses                        []uint             `json:"courses" validate:"required"`
	Keywords                       []string           `json:"keywords"`
	DeadlineAt                     string             `json:"deadlineAt" validate:"required"`
	PaymentForAccommodationID      uint               `json:"paymentForAccommodationID" validate:"required"`
	PaymentForAccommodationDetails string             `json:"paymentForAccommodationDetails"`
	FarePaymentID                  uint               `json:"farePaymentID" validate:"required"`
	FarePaymentDetails             string             `json:"farePaymentDetails"`
	Description                    VacancyDescription `json:"description"`
}

type UpdateRepository interface {
	GetVacancyByID(vacancyID uint) (*db_models.Vacancy, error)
	UpdateVacancy(vacancy *db_models.Vacancy) error
}

func Update(log *slog.Logger, repo UpdateRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.Update"
		log := log.With(slog.String("fn", fn))

		_ = log

		vacancyId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil || vacancyId == 0 {
			log.Info("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must be a valid positive integer"))
			return
		}

		var req UpdateVacancyRequest
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

		vacancy, err := repo.GetVacancyByID(uint(vacancyId))
		if err != nil {
			// FIX: Может быть не найдена запись
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		updateVacancy(vacancy, req)
		err = repo.UpdateVacancy(vacancy)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse("success"))
	}
}

func updateVacancy(vacancy *db_models.Vacancy, req UpdateVacancyRequest) {
	vacancy.CompanyName = &req.CompanyName
	vacancy.Name = req.VacancyName
	vacancy.FormatID = req.FormatID
	vacancy.DeadlineAt = req.DeadlineAt
	vacancy.PaymentForAccommodationID = req.PaymentForAccommodationID
	vacancy.PaymentForAccommodationDetails = utils.ToNullString(req.PaymentForAccommodationDetails)
	vacancy.FarePaymentID = req.FarePaymentID
	vacancy.FarePaymentDetails = utils.ToNullString(req.FarePaymentDetails)

	vacancy.Description.Workplace = utils.ToNullString(req.Description.Workplace)
	vacancy.Description.Position = utils.ToNullString(req.Description.Position)
	vacancy.Description.Salary = utils.ToNullString(req.Description.Salary)
	vacancy.Description.Requirements = utils.ToNullString(req.Description.Requirements)
	vacancy.Description.Food = utils.ToNullString(req.Description.Food)
	vacancy.Description.Conditions = utils.ToNullString(req.Description.Conditions)
	vacancy.Description.AdditionalInfo = utils.ToNullString(req.Description.AdditionalInfo)

	// Обновляем Courses (many-to-many)
	vacancy.Courses = make([]db_models.Course, 0, len(req.Courses))
	for _, course := range req.Courses {
		vacancy.Courses = append(vacancy.Courses, db_models.Course{ID: course})
	}

	// Обновляем Keywords (one-to-many)
	vacancy.Keywords = make([]db_models.VacancyKeywords, 0, len(req.Keywords))
	for _, keyword := range req.Keywords {
		vacancy.Keywords = append(vacancy.Keywords, db_models.VacancyKeywords{
			VacancyID: vacancy.ID, // ВАЖНО
			Keyword:   keyword,
		})
	}
}
