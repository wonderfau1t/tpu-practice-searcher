package vacancies

import (
	"errors"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"tpu-practice-searcher/internal/http-server/middlewares"
	"tpu-practice-searcher/internal/storage"
	"tpu-practice-searcher/internal/utils"
)

type ListResponse struct {
	TotalCount int          `json:"totalCount"`
	Vacancies  []VacancyDTO `json:"vacancies"`
}

type VacancyDTO struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	CompanyName     string `json:"companyName,omitempty"`
	CountOfReplies  *int   `json:"countOfReplies,omitempty"`
	IsCreatedByUser *bool  `json:"isCreatedByUser,omitempty"`
}

func List(log *slog.Logger, repo VacancyRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.List"
		log := log.With(slog.String("fn", fn))
		_ = log

		claims, ok := middlewares.CtxClaims(r.Context())
		if !ok {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, utils.NewErrorResponse("Failed to parse claims"))
			return
		}

		var dtos []VacancyDTO
		var err error
		switch claims.Role {
		case "student":
			dtos, err = getVacanciesForStudent(repo)
		case "moderator":
			dtos, err = getVacanciesForModerator(claims.UserID, repo)
		case "admin":
			dtos, err = getVacanciesForAdmin(claims.UserID, repo)
		case "HR", "headHR":
			dtos, err = getVacanciesForHR(claims.CompanyID, repo)
		}

		if err != nil {
			if errors.Is(err, storage.ErrRecordNotFound) {
				render.Status(r, http.StatusOK)
				render.JSON(w, r, utils.NewSuccessResponse(ListResponse{
					TotalCount: 0,
					Vacancies:  nil,
				}))
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, utils.NewErrorResponse("Internal server error"))
			return
		}

		response := ListResponse{
			TotalCount: len(dtos),
			Vacancies:  dtos,
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, utils.NewSuccessResponse(response))
	}
}

func getVacanciesForStudent(repo VacancyRepository) ([]VacancyDTO, error) {
	vacancies, err := repo.GetAllVacancies()
	if err != nil {
		return nil, err
	}

	dtos := make([]VacancyDTO, len(vacancies))
	for i, vacancy := range vacancies {
		var companyName string
		if vacancy.CompanyID == nil {
			companyName = *vacancy.CompanyName
		} else {
			companyName = vacancy.Company.Name
		}
		dtos[i] = VacancyDTO{
			ID:          vacancy.ID,
			Name:        vacancy.Name,
			CompanyName: companyName,
		}
	}
	return dtos, nil
}

func getVacanciesForModerator(moderatorID int64, repo VacancyRepository) ([]VacancyDTO, error) {
	departmentID, err := repo.GetDepartmentByModeratorID(moderatorID)
	if err != nil {
		return nil, err
	}
	vacancies, err := repo.GetVacanciesByDepartmentID(departmentID)
	if err != nil {
		return nil, err
	}

	dtos := make([]VacancyDTO, len(vacancies))
	for i, vacancy := range vacancies {
		var companyName string
		if vacancy.CompanyID == nil {
			companyName = *vacancy.CompanyName
		} else {
			companyName = vacancy.Company.Name
		}
		dtos[i] = VacancyDTO{
			ID:             vacancy.ID,
			Name:           vacancy.Name,
			CompanyName:    companyName,
			CountOfReplies: &vacancy.NumberOfResponses,
		}
		if vacancy.HrID == moderatorID {
			t := true
			dtos[i].IsCreatedByUser = &t
		} else {
			f := false
			dtos[i].IsCreatedByUser = &f
		}
	}
	return dtos, nil
}

func getVacanciesForAdmin(adminID int64, repo VacancyRepository) ([]VacancyDTO, error) {
	vacancies, err := repo.GetAllVacancies()
	if err != nil {
		return nil, err
	}

	dtos := make([]VacancyDTO, len(vacancies))
	for i, vacancy := range vacancies {
		var companyName string
		if vacancy.CompanyID == nil {
			companyName = *vacancy.CompanyName
		} else {
			companyName = vacancy.Company.Name
		}
		dtos[i] = VacancyDTO{
			ID:             vacancy.ID,
			Name:           vacancy.Name,
			CompanyName:    companyName,
			CountOfReplies: &vacancy.NumberOfResponses,
		}
		if vacancy.HrID == adminID {
			t := true
			dtos[i].IsCreatedByUser = &t
		} else {
			f := false
			dtos[i].IsCreatedByUser = &f
		}
	}
	return dtos, nil
}

func getVacanciesForHR(companyID uint, repo VacancyRepository) ([]VacancyDTO, error) {
	vacancies, err := repo.GetAllVacanciesOfCompany(companyID)
	if err != nil {
		return nil, err
	}
	t := true
	dtos := make([]VacancyDTO, len(vacancies))
	for i, vacancy := range vacancies {
		var companyName string
		if vacancy.CompanyID == nil {
			companyName = *vacancy.CompanyName
		} else {
			companyName = vacancy.Company.Name
		}
		dtos[i] = VacancyDTO{
			ID:              vacancy.ID,
			Name:            vacancy.Name,
			CompanyName:     companyName,
			CountOfReplies:  &vacancy.NumberOfResponses,
			IsCreatedByUser: &t,
		}
	}
	return dtos, nil
}
