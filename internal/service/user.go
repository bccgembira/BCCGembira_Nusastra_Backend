package service

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/entity"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/repository"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/gomail"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/helper"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/jwt"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/supabase"
	val "github.com/CRobinDev/BCCGembira_Nusastra/pkg/validator"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
)

type IUserService interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	GoogleRegister(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	GoogleLogin(ctx context.Context, req dto.GoogleLoginRequest) (dto.LoginResponse, error)
	GetUser(ctx context.Context, req dto.TokenLoginRequest) (dto.TokenLoginResponse, error)
	Update(ctx context.Context, req dto.UpdateRequest) error
	Delete(ctx context.Context, req dto.DeleteRequest) error
	FindByEmail(ctx context.Context, email string) (dto.ConvertUserEntity, error)
	SendNotification(ctx context.Context, req dto.NotificationRequest) error
	UploadProfileImage(ctx context.Context, id uuid.UUID, file *multipart.FileHeader) (string, error)
}

type userService struct {
	ur       repository.IUserRepository
	jwt      jwt.IJWT
	gomail   *gomail.Gomail
	supabase supabase.ISupabase
	logger   *log.Logger
}

func NewUserService(ur repository.IUserRepository, jwt jwt.IJWT, gomail *gomail.Gomail, supabase supabase.ISupabase, logger *log.Logger) IUserService {
	return &userService{
		ur:       ur,
		jwt:      jwt,
		gomail:   gomail,
		supabase: supabase,
		logger:   logger,
	}
}

func (us *userService) Register(ctx context.Context, req dto.RegisterRequest) error {
	if err := val.ValidateRequestRegister(req); err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		}).Error("[userService.Register] invalid request")
		return err
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		}).Error("[userService.Register] failed to hash password")
		return &response.ErrHashPassword
	}

	time := helper.GetCurrentTime()

	user := entity.User{
		ID:          uuid.New(),
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Password:    hashedPassword,
		CreatedAt:   time,
		UpdatedAt:   time,
	}

	if err := us.ur.Create(ctx, &user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			us.logger.WithFields(map[string]interface{}{
				"error": err.Error(),
			}).Error("[userService.Register] failed to register user")
			return &response.ErrUserAlreadyExists
		}
		return err
	}

	us.logger.WithFields(map[string]interface{}{
		"email": req.Email,
	}).Info("[userService.Register] user registered")

	return nil
}

func (us *userService) GoogleRegister(ctx context.Context, req dto.RegisterRequest) error {
	time := helper.GetCurrentTime()

	user := entity.User{
		ID:          uuid.New(),
		DisplayName: req.DisplayName,
		Email:       req.Email,
		CreatedAt:   time,
		UpdatedAt:   time,
	}

	if err := us.ur.Create(ctx, &user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			us.logger.WithFields(map[string]interface{}{
				"error": err.Error(),
			}).Error("[userService.GoogleRegister] failed to register user to db")
			return &response.ErrUserAlreadyExists
		}
		return err
	}

	us.logger.WithFields(map[string]interface{}{
		"email": req.Email,
	}).Info("[userService.GoogleRegister] google user registered to db")

	return nil
}

func (us *userService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := us.ur.FindByEmail(ctx, req.Email)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		}).Error("[userService.Login] failed to find user by email")
		return dto.LoginResponse{}, err
	}

	if err := helper.ComparePassword(user.Password, req.Password); err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		}).Error("[userService.Login] invalid password")
		return dto.LoginResponse{}, &response.ErrCredentialMismatch
	}

	token, err := us.jwt.CreateToken(&user)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		}).Error("[userService.Login] failed to create token")
		return dto.LoginResponse{}, &response.ErrJWTToken
	}

	resp := dto.LoginResponse{
		DisplayName: user.DisplayName,
		ID:          user.ID,
		Token:       token,
	}

	us.logger.WithFields(map[string]interface{}{
		"email": user.Email,
	}).Info("[userService.Login] user logged in")
	return resp, nil
}

func (us *userService) GoogleLogin(ctx context.Context, req dto.GoogleLoginRequest) (dto.LoginResponse, error) {
	user, err := us.ur.FindByEmail(ctx, req.Email)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		}).Error("[userService.GoogleLogin] failed to find user by email")
		return dto.LoginResponse{}, err
	}

	token, err := us.jwt.CreateToken(&user)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		}).Error("[userService.GoogleLogin] failed to create jwt token")
		return dto.LoginResponse{}, &response.ErrJWTToken
	}

	us.logger.WithFields(map[string]interface{}{
		"email": user.Email,
	}).Info("[userService.GoogleLogin] user logged in")

	return dto.LoginResponse{
		DisplayName: user.DisplayName,
		ID:          user.ID,
		Token:       token,
	}, nil
}

func (us *userService) GetUser(ctx context.Context, req dto.TokenLoginRequest) (dto.TokenLoginResponse, error) {
	user, err := us.ur.FindByID(ctx, req.ID)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    req.ID,
		}).Error("[userService.GetUser] failed to find user by ID")
		return dto.TokenLoginResponse{}, err
	}

	us.logger.WithFields(map[string]interface{}{
		"username": user.DisplayName,
	}).Info("[userService.GetUser] success get user")

	return dto.TokenLoginResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}, nil
}

func (us *userService) Update(ctx context.Context, req dto.UpdateRequest) error {
	user, err := us.ur.FindByID(ctx, req.ID)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    req.ID,
		}).Error("[userService.Update] failed to find user by ID")
		return err
	}

	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}

	if req.NewPassword != "" {
		hashedPassword, err := helper.HashPassword(req.NewPassword)
		if err != nil {
			us.logger.WithFields(map[string]interface{}{
				"error": err.Error(),
				"id":    req.ID,
			}).Error("[userService.Update] failed to hash password")
			return &response.ErrHashPassword
		}
		user.Password = hashedPassword
	}

	time := helper.GetCurrentTime()

	user.UpdatedAt = time

	rowsAffected, err := us.ur.Update(ctx, &user)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    req.ID,
		}).Error("[userService.Update] failed to update user to db")
		return err
	}

	if rowsAffected == 0 {
		us.logger.WithFields(map[string]interface{}{
			"id": req.ID,
		}).Error("[userService.Update] user not found")
		return &response.ErrUserNotFound
	}

	us.logger.WithFields(map[string]interface{}{
		"username": user.DisplayName,
	}).Info("[userService.Update] user updated")

	return nil
}

func (us *userService) Delete(ctx context.Context, req dto.DeleteRequest) error {
	rowsAffected, err := us.ur.Delete(ctx, req.ID)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    req.ID,
		}).Error("[userService.Delete] failed to delete user from db")
	}

	if rowsAffected == 0 {
		us.logger.WithFields(map[string]interface{}{
			"id": req.ID,
		}).Error("[userService.Delete] user not found")
		return &response.ErrUserNotFound
	}

	us.logger.WithFields(map[string]interface{}{
		"id": req.ID,
	}).Info("[userService.Delete] user deleted")

	return nil
}

func (us *userService) FindByEmail(ctx context.Context, email string) (dto.ConvertUserEntity, error) { 
	user, err := us.ur.FindByEmail(ctx, email)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"email": email,
		}).Error("[userService.FindByEmail] failed to find user by email")
		return dto.ConvertUserEntity{}, err
	}

	us.logger.WithFields(map[string]interface{}{
		"email": email,
	}).Info("[userService.FindByEmail] user found by email")
	return dto.ConvertUserEntity{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Image:       user.Image,
		Point:       user.Point,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (us *userService) SendNotification(ctx context.Context, req dto.NotificationRequest) error {
	user, err := us.ur.FindByID(ctx, req.ID)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    req.ID,
		}).Error("[userService.SendNotification] failed to find user by ID")
		return err
	}

	req.DisplayName = user.DisplayName
	req.Email = user.Email

	if err := us.gomail.SendNotification(req); err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    req.ID,
		}).Error("[userService.SendNotification] failed to send notification")
		return &response.ErrFailedSendNotification
	}

	return nil
}

func (us *userService) UploadProfileImage(ctx context.Context, id uuid.UUID, file *multipart.FileHeader) (string, error) {
	user, err := us.ur.FindByID(ctx, id)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("[userService.UploadProfileImage] failed to find user by ID")
		return "", err
	}

	if user.Image != "" {
		err = us.supabase.Delete(user.Image)
		if err != nil {
			us.logger.WithFields(map[string]interface{}{
				"error": err.Error(),
				"id":    id,
			}).Error("[userService.UploadProfileImage] failed to delete old profile image")
			return "", &response.ErrFailedDeleteImage
		}
	}

	url, err := us.supabase.Upload(file)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("[userService.UploadProfileImage] failed to upload profile image")
		return "", err
	}

	err = us.ur.UploadProfileImage(ctx, id, url)
	if err != nil {
		us.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		}).Error("[userService.UploadProfileImage] failed to update profile image")
		return "", err
	}

	return url, nil
}
