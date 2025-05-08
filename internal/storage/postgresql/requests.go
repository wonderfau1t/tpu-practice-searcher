package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tpu-practice-searcher/internal/http-server/handlers/get_all_hrs_of_company"
	"tpu-practice-searcher/internal/storage"
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
	if err := s.db.Debug().
		Preload("Company").
		Preload("Courses").
		Where("company_id = ?", companyID).
		Where("status_id = ?", 5).
		Find(&vacancies).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrRecordNotFound
		}
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

func (s *Storage) GetAllHrsOfCompany(companyID uint) ([]get_all_hrs_of_company.HRDTO, error) {
	var response []get_all_hrs_of_company.HRDTO
	var hrs []db_models.HrManager
	if err := s.db.Preload("User").Where("company_id = ?", companyID).Find(&hrs).Error; err != nil {
		return nil, err
	}

	for _, hr := range hrs {
		var vacancyCount int64
		if err := s.db.Model(&db_models.Vacancy{}).Where("hr_id = ?", hr.UserID).Count(&vacancyCount).Error; err != nil {
			return nil, err
		}

		response = append(response, get_all_hrs_of_company.HRDTO{
			Id:               hr.UserID,
			Username:         hr.User.Username,
			CountOfVacancies: int(vacancyCount),
		})
	}

	return response, nil
}

func (s *Storage) CreateNewHr(username string, companyID uint) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		user := db_models.User{
			Username: username,
			StatusID: constants.StatusUnderReview,
			RoleID:   constants.RoleHR,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		relation := db_models.HrManager{
			UserID:    user.ID,
			CompanyID: companyID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetCompanyInfo(companyID uint) (*db_models.Company, error) {
	var company db_models.Company
	if err := s.db.First(&company, companyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}
	return &company, nil
}

func (s *Storage) GetVacancyByID(vacancyID uint) (*db_models.Vacancy, error) {
	var vacancy db_models.Vacancy
	if err := s.db.Debug().
		Preload("Company").
		Preload("Hr").
		Preload("Format").
		Preload("Category").
		Preload("PaymentForAccommodation").
		Preload("FarePayment").
		Preload("Description").
		Preload("Courses").
		Preload("Keywords").
		Preload("Replies").
		Preload("Replies.Student").
		Where("status_id = ?", 5).
		First(&vacancy, vacancyID).Error; err != nil {
		return nil, err
	}
	return &vacancy, nil
}

func (s *Storage) GetAllVacancies() ([]db_models.Vacancy, error) {
	var vacancies []db_models.Vacancy
	if err := s.db.Debug().
		Preload("Company").
		Preload("Category").
		Preload("Courses").
		Where("status_id = ?", 5).
		Find(&vacancies).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrRecordNotFound
		}
		return nil, err
	}
	return vacancies, nil
}

func (s *Storage) AddReply(studentID int64, vacancyID uint) error {
	var vacancy db_models.Vacancy
	if err := s.db.Find(&vacancy, vacancyID).Error; err != nil {
		return err
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		reply := db_models.Reply{
			VacancyID: vacancyID,
			StudentID: studentID,
		}
		if err := tx.Create(&reply).Error; err != nil {
			return err
		}
		vacancy.NumberOfResponses = vacancy.NumberOfResponses + 1
		tx.Save(&vacancy)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetRepliedVacancies(studentID int64) ([]db_models.Vacancy, error) {
	var vacancies []db_models.Vacancy

	var replies []db_models.Reply
	if err := s.db.Preload("Vacancy.Category").Preload("Vacancy.Company").Find(&replies, "student_id = ?", studentID).Error; err != nil {
		return nil, err
	}

	for _, reply := range replies {
		vacancies = append(vacancies, reply.Vacancy)
	}
	return vacancies, nil
}

func (s *Storage) GetVacanciesBySchoolID(schoolID uint) ([]db_models.Vacancy, error) {
	var vacancies []db_models.Vacancy

	result := s.db.Distinct().
		Joins("JOIN vacancy_courses ON vacancy_courses.vacancy_id = vacancies.id").
		Joins("JOIN courses ON courses.id = vacancy_courses.course_id").
		Where("courses.school_id = ?", schoolID).
		Where("vacancies.status_id = ?", 5).
		Preload("Company").
		Preload("Status").
		Preload("Format").
		Preload("Category").
		Preload("Hr").
		Preload("PaymentForAccommodation").
		Preload("FarePayment").
		Preload("Courses").
		Preload("Description").
		Preload("Keywords").
		Find(&vacancies)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, storage.ErrRecordNotFound
		}
		return nil, result.Error
	}

	return vacancies, nil
}

func (s *Storage) GetSchoolByModeratorID(moderatorID int64) (uint, error) {
	var moderator db_models.Moderator
	if err := s.db.First(&moderator, moderatorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, storage.ErrRecordNotFound
		}
		return 0, err
	}
	return moderator.SchoolID, nil
}

func (s *Storage) DeleteReply(studentID int64, vacancyID uint) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Удаление отклика
		if err := tx.Where("vacancy_id = ? AND student_id = ?", vacancyID, studentID).
			Delete(&db_models.Reply{}).Error; err != nil {
			return err
		}
		// Уменьшить счетчик
		if err := tx.Model(&db_models.Vacancy{}).
			Where("id = ?", vacancyID).
			UpdateColumn("number_of_responses", gorm.Expr("number_of_responses - ?", 1)).
			Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *Storage) FilterVacancies(categoryID *uint, courseIDs []uint) ([]db_models.Vacancy, error) {
	var vacancies []db_models.Vacancy
	query := s.db.Model(&db_models.Vacancy{})

	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if len(courseIDs) > 0 {
		query = query.Joins("JOIN vacancy_courses vc ON vc.vacancy_id = vacancies.id").
			Where("vc.course_id IN ?", courseIDs).
			Group("vacancies.id")
	}

	if err := query.Preload("Category").
		Preload("Company").
		Find(&vacancies).Error; err != nil {
		return nil, err
	}

	return vacancies, nil
}

func (s *Storage) IsReplied(studentID int64, vacancyID uint) (*bool, error) {
	var count int64
	err := s.db.Model(&db_models.Reply{}).
		Where("student_id = ? AND vacancy_id = ?", studentID, vacancyID).
		Count(&count).Error
	if err != nil {
		return nil, err
	}
	result := count > 0
	return &result, nil
}

func (s *Storage) SearchVacancies(searchQuery string) ([]db_models.Vacancy, error) {
	var vacancies []db_models.Vacancy
	query := s.db.Model(&db_models.Vacancy{})

	searchPattern := "%" + searchQuery + "%"
	query = query.Joins("LEFT JOIN vacancy_keywords vk ON vk.vacancy_id = vacancies.id").
		Where("vacancies.name ILIKE ? OR vk.keyword ILIKE ?", searchPattern, searchPattern).
		Group("vacancies.id")

	if err := query.Preload("Category").
		Preload("Company").
		//Preload("Courses").
		//Preload("Description").
		//Preload("Keywords").
		Find(&vacancies).Error; err != nil {
		return nil, err
	}

	return vacancies, nil
}

func (s *Storage) FindOrRegisterByUsername(userID int64, username string) (*db_models.User, bool, error) {
	var user db_models.User

	// Сначала ищем по userID
	err := s.db.Preload("Role").Where("id = ?", userID).First(&user).Error
	if err == nil {
		return &user, true, nil // Пользователь уже полностью зарегистрирован
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err // Возвращаем ошибку, если это не "не найдено"
	}

	err = s.db.Preload("Role").Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, storage.ErrRecordNotFound // Пользователь вообще не существует
	}
	if err != nil {
		return nil, false, err // Другая ошибка
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {

		var hr db_models.HrManager
		err := tx.Find(&hr, user.ID).Error
		if err != nil {
			return err
		}
		companyID := hr.CompanyID

		err = tx.Delete(&hr).Error
		if err != nil {
			return err
		}

		err = tx.Model(&user).Updates(map[string]interface{}{
			"id":        userID,
			"status_id": 5,
		}).Error
		if err != nil {
			return err
		}

		err = tx.Create(&db_models.HrManager{UserID: userID, CompanyID: companyID}).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, false, err
	}
	user.ID = userID
	return &user, true, nil
}

func (s *Storage) ApprovePhoneNumber(userID int64, phoneNumber string) error {
	err := s.db.Model(&db_models.User{}).
		Where("id = ?", userID).
		UpdateColumn("phone_number", phoneNumber).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) HideVacancyByID(vacancyID uint) error {
	err := s.db.Model(&db_models.Vacancy{}).
		Where("id = ?", vacancyID).
		UpdateColumn("status_id", 4).Error
	return err
}

func (s *Storage) UpdateCompanyInfo(companyID uint, name string, description string, link string) error {
	err := s.db.Model(&db_models.Company{}).
		Where("id = ?", companyID).
		Updates(map[string]interface{}{
			"name":        name,
			"description": description,
			"link":        link,
		}).Error
	return err
}
