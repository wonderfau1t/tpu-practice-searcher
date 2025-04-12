package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tpu-practice-searcher/internal/storage/models/db_models"
	"tpu-practice-searcher/internal/utils/constants"
)

func (s *Storage) GetUserByID(userID int64) (*db_models.User, error) {
	const fn = "storage.postgresql.GetUserByID"

	var user db_models.User
	err := s.db.Preload("Role").First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &user, nil
}

func (s *Storage) GetAllCategories() ([]db_models.Category, error) {
	const fn = "storage.postgresql.GetAllCategories"

	var categories []db_models.Category
	err := s.db.Debug().Find(&categories).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return categories, nil
}

func (s *Storage) GetAllFormats() ([]db_models.Format, error) {
	const fn = "storage.postgresql.GetAllFormats"

	var format []db_models.Format
	err := s.db.Debug().Find(&format).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return format, nil
}

func (s *Storage) GetAllFarePaymentMethods() ([]db_models.FarePayment, error) {
	const fn = "storage.postgresql.GetAllFarePaymentMethod"

	var methods []db_models.FarePayment
	err := s.db.Debug().Find(&methods).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return methods, nil
}

func (s *Storage) GetAllAccommodationPaymentMethods() ([]db_models.PaymentForAccommodation, error) {
	const fn = "storage.postgresql.GetAllAccommodationPaymentMethod"

	var methods []db_models.PaymentForAccommodation
	err := s.db.Debug().Find(&methods).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return methods, nil
}

func (s *Storage) GetAllCourses() ([]db_models.Course, error) {
	const fn = "storage.postgresql.GetAllAccommodationPaymentMethod"

	var courses []db_models.Course
	err := s.db.Debug().Find(&courses).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return courses, nil
}

func (s *Storage) IsUserExist(userId int64) (bool, error) {
	const fn = "storage.postgresql.IsEmailExist"

	err := s.db.Where("id = ?", userId).First(&db_models.User{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return true, nil
}

func (s *Storage) CreateNewUser(userId int64, username string, roleId uint) error {
	const fn = "storage.postgresql.CreateUserAuth"

	user := db_models.User{
		ID:       userId,
		Username: username,
		StatusID: constants.StatusDefault,
		RoleID:   roleId,
	}
	err := s.db.Create(&user).Error
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) CreateNewCompany(userId int64, username string, companyName string, companyDescription sql.NullString, companyLink sql.NullString) error {
	const fn = "storage.postgresql.CreateNewCompany"

	err := s.db.Transaction(func(tx *gorm.DB) error {
		hr := db_models.User{
			ID:       userId,
			Username: username,
			StatusID: constants.StatusDefault,
			RoleID:   constants.RoleHeadHR,
		}
		if err := tx.Create(&hr).Error; err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		company := db_models.Company{
			Name:        companyName,
			Description: companyDescription,
			Link:        companyLink,
			StatusID:    constants.StatusDefault,
			HeadHrID:    hr.ID,
		}
		if err := tx.Create(&company).Error; err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		relation := db_models.HrManager{
			UserID:    hr.ID,
			CompanyID: company.ID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateNewVacancy(vacancy *db_models.Vacancy) error {
	if err := s.db.Create(vacancy).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetCompanyByHrID(hrID int64) (*db_models.HrManager, error) {
	const fn = "storage.postgresql.GetCompanyByHrID"

	var company db_models.HrManager
	if err := s.db.First(&company, "user_id = ?", hrID).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (s *Storage) IsCompanyExist(companyID uint) (bool, error) {
	const fn = "storage.postgresql.IsCompanyExist"

	err := s.db.Where("id = ?", companyID).First(&db_models.Company{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return true, nil
}

func (s *Storage) GetAllVacanciesOfCompany(companyID uint) ([]db_models.Vacancy, error) {
	//const fn = "storage.postg"

	var vacancies []db_models.Vacancy
	if err := s.db.Debug().Preload("Courses").Preload("Category").Where("company_id = ?", companyID).Find(&vacancies).Error; err != nil {
		return nil, err
	}
	return vacancies, nil
}

func (s *Storage) GetCompanyIDByHRID(hrID int64) (uint, error) {

	var companyHr db_models.HrManager
	if err := s.db.First(&companyHr, "user_id = ?", hrID).Error; err != nil {
		return 0, err
	}
	return companyHr.CompanyID, nil
}
