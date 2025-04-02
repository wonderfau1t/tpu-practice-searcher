package postgresql

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tpu-practice-searcher/internal/storage/models"
)

func (s *Storage) GetUserByID(userID int64) (*models.User, error) {
	const fn = "storage.postgresql.GetUserByID"

	var user models.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &user, nil
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
		RoleID:   roleId,
	}
	err := s.db.Create(&user).Error
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
