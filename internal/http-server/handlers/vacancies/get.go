package vacancies

import (
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils"
)

type GetStorage interface {
	GetAllVacancies() ([]db_models.Vacancy, error)
	GetAllVacanciesOfCompany(companyID uint) ([]db_models.Vacancy, error)
	GetSchoolByModeratorID(moderatorID int64) (uint, error)
	GetVacanciesBySchoolID(schoolID uint) ([]db_models.Vacancy, error)
}

func GetVacancies(log *slog.Logger, db GetStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.vacancies.GetVacancies"
		log := log.With(slog.String("fn", fn))

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			log.Info("not valid accessToken")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("failed to parse claims"))
			return
		}

		var dtos []GetVacancyDTO
		var err error
		switch claims.Role {
		case "student":
			dtos, err = getAllVacanciesForStudent(db)
		case "HR", "headHR":
			dtos, err = getAllVacanciesForHr(claims.CompanyID, db)
		case "moderator":
			dtos, err = getAllVacanciesForModerator(claims.UserID, db)
		}

		if err != nil {
			if errors.Is(err, storage.ErrRecordNotFound) {
				render.Status(r, http.StatusOK)
				render.JSON(w, r, utils.NewSuccessResponse(GetResponse{
					TotalCount: 0,
					Vacancies:  nil,
				}))
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(GetResponse{
			TotalCount: len(dtos),
			Vacancies:  dtos,
		}))
	}
}

func getAllVacanciesForStudent(db GetStorage) ([]GetVacancyDTO, error) {
	vacancies, err := db.GetAllVacancies()
	if err != nil {
		if errors.Is(err, storage.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	dtos := make([]GetVacancyDTO, len(vacancies))
	for i, vacancy := range vacancies {
		dtos[i] = GetVacancyDTO{
			ID:          vacancy.ID,
			Name:        vacancy.Name,
			CompanyName: vacancy.Company.Name,
			Category:    vacancy.Category.Name,
		}
	}
	return dtos, err
}

func getAllVacanciesForHr(companyID uint, db GetStorage) ([]GetVacancyDTO, error) {
	vacancies, err := db.GetAllVacanciesOfCompany(companyID)
	if err != nil {
		return nil, err
	}

	dtos := make([]GetVacancyDTO, len(vacancies))
	for i, vacancy := range vacancies {
		print(vacancy.NumberOfResponses)
		dtos[i] = GetVacancyDTO{
			ID:             vacancy.ID,
			Name:           vacancy.Name,
			Category:       vacancy.Category.Name,
			CountOfReplies: &vacancy.NumberOfResponses,
		}
	}
	return dtos, err
}

func getAllVacanciesForModerator(moderatorID int64, db GetStorage) ([]GetVacancyDTO, error) {
	schoolID, err := db.GetSchoolByModeratorID(moderatorID)
	if err != nil {
		return nil, err
	}

	vacancies, err := db.GetVacanciesBySchoolID(schoolID)
	if err != nil {
		return nil, err
	}

	dtos := make([]GetVacancyDTO, len(vacancies))
	for i, vacancy := range vacancies {
		dtos[i] = GetVacancyDTO{
			ID:             vacancy.ID,
			Name:           vacancy.Name,
			CompanyName:    vacancy.Company.Name,
			Category:       vacancy.Category.Name,
			CountOfReplies: &vacancy.NumberOfResponses,
		}
	}
	return dtos, nil
}
