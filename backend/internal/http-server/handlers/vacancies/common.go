package vacancies

import (
	"errors"
	"tpu-practice-searcher/internal/storage"
	"tpu-practice-searcher/internal/storage/models/db_models"
)

type GetVacancyDTO struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	CompanyName    string `json:"companyName,omitempty"`
	CountOfReplies *int   `json:"countOfReplies,omitempty"`
}

type GetResponse struct {
	TotalCount int             `json:"totalCount"`
	Vacancies  []GetVacancyDTO `json:"vacancies"`
}

type DetailsVacancyDTO struct {
	Id                             uint                  `json:"id"`
	IsCreatedByUser                *bool                 `json:"isCreatedByUser,omitempty"`
	VacancyName                    string                `json:"vacancyName"`
	HasCompanyProfile              *bool                 `json:"hasCompanyProfile,omitempty"`
	IsReplied                      *bool                 `json:"isReplied,omitempty"`
	CompanyID                      *uint                 `json:"companyID"`
	CompanyName                    string                `json:"companyName"`
	Format                         string                `json:"format"`
	FormatID                       uint                  `json:"formatID"`
	Courses                        []interface{}         `json:"courses"`
	Keywords                       []string              `json:"keywords"`
	DeadlineAt                     string                `json:"deadlineAt"`
	PaymentForAccommodationID      uint                  `json:"paymentForAccommodationID"`
	PaymentForAccommodation        string                `json:"paymentForAccommodation"`
	PaymentForAccommodationDetails string                `json:"paymentForAccommodationDetails,omitempty"`
	FarePaymentID                  uint                  `json:"farePaymentID"`
	FarePayment                    string                `json:"farePayment"`
	FarePaymentDetails             string                `json:"farePaymentDetails,omitempty"`
	Description                    DetailsDescriptionDTO `json:"description"`
	HrInfo                         *DetailsUserDTO       `json:"hrInfo,omitempty"`
	RepliedStudents                []DetailsUserDTO      `json:"repliedStudents,omitempty"`
}

type DetailsDescriptionDTO struct {
	Workplace      string `json:"workplace"`
	Position       string `json:"position"`
	Salary         string `json:"salary"`
	Requirements   string `json:"requirements"`
	Food           string `json:"food"`
	Conditions     string `json:"conditions"`
	AdditionalInfo string `json:"additionalInfo"`
}

type DetailsUserDTO struct {
	ID          int64  `json:"id,omitempty"`
	Username    string `json:"username,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type DetailsStorage interface {
	GetVacancyByID(vacancyID uint) (*db_models.Vacancy, error)
	IsReplied(studentID int64, vacancyID uint) (*bool, error)
}

func GetAllVacanciesForStudent(db GetStorage) ([]GetVacancyDTO, error) {
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
		}
	}
	return dtos, err
}

func GetAllVacanciesForHr(companyID uint, db GetStorage) ([]GetVacancyDTO, error) {
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
			CountOfReplies: &vacancy.NumberOfResponses,
		}
	}
	return dtos, err
}

func GetAllVacanciesForModerator(moderatorID int64, db GetStorage) ([]GetVacancyDTO, error) {
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
			CountOfReplies: &vacancy.NumberOfResponses,
		}
	}
	return dtos, nil
}
