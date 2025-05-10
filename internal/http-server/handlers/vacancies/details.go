package vacancies

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type DetailsResponse struct {
	VacancyInfo DetailsVacancyDTO `json:"vacancyInfo"`
}

func GetVacancyDetails(log *slog.Logger, db DetailsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.vacancies.GetVacancyDetails"
		log := log.With(slog.String("fn", fn))

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			log.Info("not valid accessToken")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("failed to parse claims"))
			return
		}

		vacancyId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil || vacancyId == 0 {
			log.Info("invalid id")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must 1be a valid positive integer"))
			return
		}

		vacancy, err := db.GetVacancyByID(uint(vacancyId))
		if err != nil {
			if errors.Is(err, storage.ErrRecordNotFound) {
				log.Info("vacancy not found")
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, utils.NewSuccessResponse("vacancy not found"))
				return
			}
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		var dto DetailsVacancyDTO

		switch claims.Role {
		case "student":
			dto = toVacancyDTO(vacancy)
			dto.IsReplied, _ = db.IsReplied(claims.UserID, uint(vacancyId))
		case "headHR", "HR":
			dto = toVacancyDTO(vacancy)
		case "moderator", "admin":
			dto = toVacancyDTOFull(vacancy)
			t := new(bool)
			f := new(bool)

			*t = true
			*f = false
			if claims.UserID == dto.HrInfo.ID {
				dto.IsCreatedByUser = t
			} else {
				dto.IsCreatedByUser = f
			}
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(DetailsResponse{VacancyInfo: dto}))
	}
}

func toVacancyDTO(vacancy *db_models.Vacancy) DetailsVacancyDTO {
	dto := DetailsVacancyDTO{
		Id:                        vacancy.ID,
		VacancyName:               vacancy.Name,
		CompanyID:                 vacancy.CompanyID,
		FormatID:                  vacancy.FormatID,
		Format:                    vacancy.Format.Name,
		DeadlineAt:                vacancy.DeadlineAt,
		PaymentForAccommodationID: vacancy.PaymentForAccommodationID,
		PaymentForAccommodation:   vacancy.PaymentForAccommodation.Name,
		FarePaymentID:             vacancy.FarePaymentID,
		FarePayment:               vacancy.FarePayment.Name,
		Description: DetailsDescriptionDTO{
			Workplace:      vacancy.Description.Workplace.String,
			Position:       vacancy.Description.Position.String,
			Salary:         vacancy.Description.Salary.String,
			Requirements:   vacancy.Description.Requirements.String,
			Food:           vacancy.Description.Food.String,
			Conditions:     vacancy.Description.Conditions.String,
			AdditionalInfo: vacancy.Description.AdditionalInfo.String,
		},
	}
	if dto.FarePaymentID == 3 {
		dto.FarePaymentDetails = vacancy.FarePaymentDetails.String
	}
	if dto.PaymentForAccommodationID == 3 {
		dto.PaymentForAccommodationDetails = vacancy.PaymentForAccommodationDetails.String
	}
	if vacancy.CompanyID != nil {
		dto.CompanyName = vacancy.Company.Name
	} else {
		dto.CompanyName = *vacancy.CompanyName
	}
	for _, course := range vacancy.Courses {
		courseLocal := map[string]string{
			"courseId": strconv.Itoa(int(course.ID)),
			"name":     course.Name,
		}
		dto.Courses = append(dto.Courses, courseLocal)
	}
	for _, keyword := range vacancy.Keywords {
		dto.Keywords = append(dto.Keywords, keyword.Keyword)
	}
	return dto
}

func toVacancyDTOFull(vacancy *db_models.Vacancy) DetailsVacancyDTO {
	dto := DetailsVacancyDTO{
		Id:                        vacancy.ID,
		VacancyName:               vacancy.Name,
		CompanyID:                 vacancy.CompanyID,
		FormatID:                  vacancy.FormatID,
		Format:                    vacancy.Format.Name,
		DeadlineAt:                vacancy.DeadlineAt,
		PaymentForAccommodationID: vacancy.PaymentForAccommodationID,
		PaymentForAccommodation:   vacancy.PaymentForAccommodation.Name,
		FarePaymentID:             vacancy.FarePaymentID,
		FarePayment:               vacancy.FarePayment.Name,
		Description: DetailsDescriptionDTO{
			Workplace:      vacancy.Description.Workplace.String,
			Position:       vacancy.Description.Position.String,
			Salary:         vacancy.Description.Salary.String,
			Requirements:   vacancy.Description.Requirements.String,
			Food:           vacancy.Description.Food.String,
			Conditions:     vacancy.Description.Conditions.String,
			AdditionalInfo: vacancy.Description.AdditionalInfo.String,
		},
		HrInfo: &DetailsUserDTO{
			ID:          vacancy.Hr.ID,
			Username:    vacancy.Hr.Username,
			PhoneNumber: vacancy.Hr.PhoneNumber.String,
		},
	}
	if vacancy.CompanyID != nil {
		dto.CompanyName = vacancy.Company.Name
	} else {
		dto.CompanyName = *vacancy.CompanyName
	}
	if dto.FarePaymentID == 3 {
		dto.FarePaymentDetails = vacancy.FarePaymentDetails.String
	}
	if dto.PaymentForAccommodationID == 3 {
		dto.PaymentForAccommodationDetails = vacancy.PaymentForAccommodationDetails.String
	}
	for _, course := range vacancy.Courses {
		courseLocal := map[string]string{
			"courseId": strconv.Itoa(int(course.ID)),
			"name":     course.Name,
		}
		dto.Courses = append(dto.Courses, courseLocal)
	}
	for _, keyword := range vacancy.Keywords {
		dto.Keywords = append(dto.Keywords, keyword.Keyword)
	}

	for _, reply := range vacancy.Replies {
		dto.RepliedStudents = append(dto.RepliedStudents, DetailsUserDTO{
			ID:          reply.Student.ID,
			Username:    reply.Student.Username,
			PhoneNumber: reply.Student.PhoneNumber.String,
		})
	}
	return dto
}
