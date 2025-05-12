package companies

import (
	"database/sql"
	"tpu-practice-searcher/internal/storage/models/db_models"
)

type CompanyRepository interface {
	CreateNewCompany(userId int64, username string, companyName string, companyDescription sql.NullString, companyLink sql.NullString) error
	GetAllHrsOfCompany(companyID uint) ([]HRDTO, error)
	CreateNewHr(username string, companyID uint) error
	GetCompanyInfo(companyID uint) (*db_models.Company, error)
}

type HRDTO struct {
	Id               int64  `json:"id"`
	Username         string `json:"username"`
	CountOfVacancies int    `json:"countOfVacancies"`
}
