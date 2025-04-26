package repository

import (
	"context"
	"errors"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/entity"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/helper"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	Update(ctx context.Context, user *entity.User) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) (int64, error)
	UploadProfileImage(ctx context.Context, id uuid.UUID, url string) error
}

type userRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewUserRepository(db *gorm.DB, logger *logrus.Logger) IUserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (ur *userRepository) Create(ctx context.Context, user *entity.User) error {
	tr := ur.db.WithContext(ctx).Begin()

	if err := tr.Create(&user).Error; err != nil {
		tr.Rollback()

		ur.logger.Errorf("failed to create user: %v", err)
		return err
	}

	if err := tr.Commit().Error; err != nil {
		tr.Rollback()
		ur.logger.Errorf("failed to commit user creation: %v", err)
		return err
	}

	return nil
}

func (ur *userRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	var user entity.User
	err := ur.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ur.logger.Warnf("user with id %s not found", id)
		return entity.User{}, &response.ErrUserNotFound
	}

	return user, err
}

func (ur *userRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	err := ur.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.Warnf("user with email %s not found", email)
			return entity.User{}, &response.ErrUserNotFound
		}
		ur.logger.Errorf("failed to find user by email %s: %v", email, err)
		return entity.User{}, err
	}
	return user, nil
}

func (ur *userRepository) Update(ctx context.Context, user *entity.User) (int64, error) {
	updatedAt := helper.GetCurrentTime()

	result := ur.db.WithContext(ctx).Model(&user).Where("id", user.ID).Updates(map[string]interface{}{
		"display_name":    user.DisplayName,
		"password":    user.Password,
		"updated_at":  updatedAt,
		"image":       user.Image,
	})

	return result.RowsAffected, result.Error
}

func (ur *userRepository) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	var user entity.User
	result := ur.db.Delete(&user, id)
	return result.RowsAffected, result.Error
}

func (ur *userRepository) UploadProfileImage(ctx context.Context, id uuid.UUID, url string) error {
	updatedAt := helper.GetCurrentTime()
	var user entity.User
	if err := ur.db.WithContext(ctx).Model(&user).Where("id = ?", id).Updates(map[string]interface{}{
		"image":      url,
		"updated_at": updatedAt,
	}).Error; err != nil {
		return err
	}

	return nil
}
