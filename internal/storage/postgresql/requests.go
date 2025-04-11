package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tpu-practice-searcher/internal/storage/models"
	"tpu-practice-searcher/internal/utils/constants"
)

func (s *Storage) GetUserByID(userID int64) (*models.User, error) {
	const fn = "storage.postgresql.GetUserByID"

	var user models.User
	err := s.db.Preload("Role").First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &user, nil
}

func (s *Storage) GetAllCategories() ([]models.Category, error) {
	const fn = "storage.postgresql.GetAllCategories"

	var categories []models.Category
	err := s.db.Debug().Find(&categories).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return categories, nil
}

func (s *Storage) GetAllFormats() ([]models.Format, error) {
	const fn = "storage.postgresql.GetAllFormats"

	var format []models.Format
	err := s.db.Debug().Find(&format).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return format, nil
}

func (s *Storage) GetAllFarePaymentMethods() ([]models.FarePayment, error) {
	const fn = "storage.postgresql.GetAllFarePaymentMethod"

	var methods []models.FarePayment
	err := s.db.Debug().Find(&methods).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return methods, nil
}

func (s *Storage) GetAllAccommodationPaymentMethods() ([]models.PaymentForAccommodation, error) {
	const fn = "storage.postgresql.GetAllAccommodationPaymentMethod"

	var methods []models.PaymentForAccommodation
	err := s.db.Debug().Find(&methods).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return methods, nil
}

func (s *Storage) GetAllCourses() ([]models.Course, error) {
	const fn = "storage.postgresql.GetAllAccommodationPaymentMethod"

	var courses []models.Course
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

	err := s.db.Where("id = ?", userId).First(&models.User{}).Error
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

	user := models.User{
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
		hr := models.User{
			ID:       userId,
			Username: username,
			StatusID: constants.StatusDefault,
			RoleID:   constants.RoleHeadHR,
		}
		if err := tx.Create(&hr).Error; err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		company := models.Company{
			Name:        companyName,
			Description: companyDescription,
			Link:        companyLink,
			StatusID:    constants.StatusDefault,
			HeadHrID:    hr.ID,
		}
		if err := tx.Create(&company).Error; err != nil {
			return fmt.Errorf("%s: %w", fn, err)
		}

		relation := models.HrManager{
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

func (s *Storage) CreateNewVacancy(vacancy *models.Vacancy) error {
	if err := s.db.Create(vacancy).Error; err != nil {
		return err
	}
	return nil
}
