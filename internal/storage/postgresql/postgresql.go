package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"tpu-practice-searcher/config/initial_data"
	"tpu-practice-searcher/internal/config"
	"tpu-practice-searcher/internal/storage/models/db_models"
)

type Storage struct {
	db *gorm.DB
}

func SetupStorage(storageConf config.Storage) (*Storage, error) {
	const fn = "storage.postrgresql.SetupStorage"

	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=%d", storageConf.User, storageConf.Password, storageConf.Dbname, storageConf.Port)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	err = db.AutoMigrate(&db_models.Role{}, &db_models.User{}, &db_models.Status{}, &db_models.Company{}, &db_models.School{}, &db_models.Course{}, &db_models.Moderator{}, &db_models.HrManager{}, &db_models.Format{}, &db_models.Category{}, &db_models.PaymentForAccommodation{}, &db_models.FarePayment{}, &db_models.Vacancy{}, &db_models.VacancyDescription{}, &db_models.VacancyKeywords{})
	if err != nil {
		log.Fatalf("failed to apply migrations")
	}

	err = initializeRoles(db)
	if err != nil {
		log.Fatalf("failed to set default roles")
	}

	err = initializeStatuses(db)
	if err != nil {
		log.Fatalf("failed to set default statuses")
	}

	err = initializeFormats(db)
	if err != nil {
		log.Fatalf("failed to set default formats")
	}

	err = initializeFarePayments(db)
	if err != nil {
		log.Fatalf("failed to set default fare payments")
	}

	err = initializePaymentForAccommodation(db)
	if err != nil {
		log.Fatalf("failed to set default PaymentForAccommodation")
	}

	err = initializeCategories(db)
	if err != nil {
		log.Fatalf("failed to set default Categories")
	}

	err = initializeSchools(db)
	if err != nil {
		log.Fatalf("failed to set default Schools")
	}

	err = initializeCourses(db)
	if err != nil {
		log.Fatalf("failed to set default Courses")
	}

	return &Storage{db: db}, nil
}

func initializeRoles(db *gorm.DB) error {
	const fn = "storage.postgresql.initializeRoles"
	for _, role := range initial_data.DefaultRoles {
		if err := db.FirstOrCreate(&role, role).Error; err != nil {
			return fmt.Errorf("%s: failed to ensure role %s: %w", fn, role.Name, err)
		}
	}
	return nil
}

func initializeStatuses(db *gorm.DB) error {
	const fn = "storage.postgresql.initializeStatuses"
	for _, status := range initial_data.DefaultStatuses {
		if err := db.FirstOrCreate(&status, status).Error; err != nil {
			return fmt.Errorf("%s: failed to add status: %s: %w", fn, status.Name, err)
		}
	}
	return nil
}

func initializeFormats(db *gorm.DB) error {
	const fn = "storage.postgresql.initializeSchools"
	for _, format := range initial_data.DefaultFormats {
		if err := db.FirstOrCreate(&format, format).Error; err != nil {
			return fmt.Errorf("%s: failed to add status: %s: %w", fn, format.Name, err)
		}
	}
	return nil
}

func initializeFarePayments(db *gorm.DB) error {
	const fn = "storage.postgresql.initializeSchools"
	for _, option := range initial_data.DefaultFarePayments {
		if err := db.FirstOrCreate(&option, option).Error; err != nil {
			return fmt.Errorf("%s: failed to add status: %s: %w", fn, option.Name, err)
		}
	}
	return nil
}

func initializePaymentForAccommodation(db *gorm.DB) error {
	const fn = "storage.postgresql.initializeSchools"
	for _, option := range initial_data.DefaultPaymentForAccommodation {
		if err := db.FirstOrCreate(&option, option).Error; err != nil {
			return fmt.Errorf("%s: failed to add status: %s: %w", fn, option.Name, err)
		}
	}
	return nil
}

func initializeCategories(db *gorm.DB) error {
	const fn = "storage.postgresql.initializeSchools"
	for _, category := range initial_data.DefaultCategories {
		if err := db.FirstOrCreate(&category, category).Error; err != nil {
			return fmt.Errorf("%s: failed to add status: %s: %w", fn, category.Name, err)
		}
	}
	return nil
}

func initializeSchools(db *gorm.DB) error {
	const fn = "storage.postgresql.initializeSchools"
	for _, school := range initial_data.DefaultSchools {
		if err := db.FirstOrCreate(&school, school).Error; err != nil {
			return fmt.Errorf("%s: failed to add status: %s: %w", fn, school.Name, err)
		}
	}
	return nil
}

func initializeCourses(db *gorm.DB) error {
	const fn = "storage.postgresql.initializeSchools"
	for _, course := range initial_data.DefaultCourses {
		if err := db.FirstOrCreate(&course, course).Error; err != nil {
			return fmt.Errorf("%s: failed to add status: %s: %w", fn, course.Name, err)
		}
	}
	return nil
}
