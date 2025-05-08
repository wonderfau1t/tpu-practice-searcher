package get_vacancy_details

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type VacancyDTO struct {
	// FIX: Информация о компании нужна еще
	Id                      uint           `json:"id"`
	Name                    string         `json:"name"`
	CompanyID               uint           `json:"companyID"`
	CompanyName             string         `json:"companyName"`
	Format                  string         `json:"format"`
	Category                string         `json:"category"`
	Courses                 []string       `json:"courses"`
	Keywords                []string       `json:"keywords"`
	DeadlineAt              string         `json:"deadlineAt"`
	PaymentForAccommodation string         `json:"paymentForAccommodation"`
	FarePayment             string         `json:"farePayment"`
	Description             DescriptionDTO `json:"description"`
}

type DescriptionDTO struct {
	Workplace      string `json:"workplace"`
	Position       string `json:"position"`
	Salary         string `json:"salary"`
	Requirements   string `json:"requirements"`
	Food           string `json:"food"`
	Conditions     string `json:"conditions"`
	AdditionalInfo string `json:"additionalInfo"`
}

type Response struct {
	VacancyInfo VacancyDTO `json:"vacancyInfo"`
}

type Storage interface {
	GetVacancyByID(vacancyID uint) (*db_models.Vacancy, error)
}

func New(log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vacancyIdStr := chi.URLParam(r, "id")

		if vacancyIdStr == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must not be empty"))
			return
		}
		vacancyId, err := strconv.ParseUint(vacancyIdStr, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, utils.NewErrorResponse("id must be type of uint"))
			return
		}

		vacancy, err := db.GetVacancyByID(uint(vacancyId))
		if err != nil {
			// FIX: Может быть не найдена запись
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		dto := VacancyDTO{
			Id:                      vacancy.ID,
			Name:                    vacancy.Name,
			CompanyID:               vacancy.CompanyID,
			CompanyName:             vacancy.Company.Name,
			Format:                  vacancy.Format.Name,
			DeadlineAt:              vacancy.DeadlineAt,
			PaymentForAccommodation: vacancy.PaymentForAccommodation.Name,
			FarePayment:             vacancy.FarePayment.Name,
			Description: DescriptionDTO{
				Workplace:      vacancy.Description.Workplace.String,
				Position:       vacancy.Description.Position.String,
				Salary:         vacancy.Description.Salary.String,
				Requirements:   vacancy.Description.Requirements.String,
				Food:           vacancy.Description.Food.String,
				Conditions:     vacancy.Description.Conditions.String,
				AdditionalInfo: vacancy.Description.AdditionalInfo.String,
			},
		}

		for _, course := range vacancy.Courses {
			dto.Courses = append(dto.Courses, course.Name)
		}
		for _, keyword := range vacancy.Keywords {
			dto.Keywords = append(dto.Keywords, keyword.Keyword)
		}

		response := Response{VacancyInfo: dto}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(response))
	}
}
