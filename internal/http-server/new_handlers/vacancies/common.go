package vacancies

import "tpu-practice-searcher/internal/storage/models/db_models"

type VacancyRepository interface {
	GetUserByID(userID int64) (*db_models.User, error)
	CreateVacancy(vacancy *db_models.Vacancy) error
	GetCompanyByHrID(hrID int64) (*db_models.HrManager, error)
	HideVacancyByID(vacancyID uint) error
	GetAllVacancies() ([]db_models.Vacancy, error)
	GetVacanciesByDepartmentID(departmentID uint) ([]db_models.Vacancy, error)
	GetDepartmentByModeratorID(moderatorID int64) (uint, error)
	GetVacancyByID(vacancyID uint) (*db_models.Vacancy, error)
	IsReplied(studentID int64, vacancyID uint) (*bool, error)
	SearchVacancies(searchQuery string) ([]db_models.Vacancy, error)
	FilterVacancies(courseIDs []uint) ([]db_models.Vacancy, error)
}

type CreateVacancyRequest struct {
	CompanyName                    string             `json:"companyName,omitempty" validate:"required"`
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

type VacancyDescription struct {
	Workplace      string `json:"workplace"`
	Position       string `json:"position"`
	Salary         string `json:"salary"`
	Requirements   string `json:"requirements"`
	Food           string `json:"food"`
	Conditions     string `json:"conditions"`
	AdditionalInfo string `json:"additionalInfo"`
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
