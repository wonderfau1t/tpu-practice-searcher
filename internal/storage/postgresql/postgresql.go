package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"tpu-practice-searcher/internal/config"
	"tpu-practice-searcher/internal/storage/models"
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

	err = db.AutoMigrate(&models.Role{}, &models.User{})
	if err != nil {
		log.Fatalf("failed to apply migrations")
	}

	err = initializeRoles(db)
	if err != nil {
		log.Fatalf("failed to set default roles")
	}

	return &Storage{db: db}, nil
}

func initializeRoles(db *gorm.DB) error {
	const fn = "storage.postgresql.initializeRoles"

	defaultRoles := []models.Role{
		{
			ID:   1,
			Name: "admin",
		},
		{
			ID:   2,
			Name: "moderator",
		},
		{
			ID:   3,
			Name: "hr",
		},
		{
			ID:   4,
			Name: "student",
		},
	}

	for _, role := range defaultRoles {
		if err := db.FirstOrCreate(&role, role).Error; err != nil {
			return fmt.Errorf("%s: failed to ensure role %s: %w", fn, role.Name, err)
		}
	}

	return nil
}
