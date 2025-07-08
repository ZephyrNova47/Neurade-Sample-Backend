package service

import (
	"be/neurade/v2/internal/entity"
	"be/neurade/v2/internal/model"
	"be/neurade/v2/internal/model/converter"
	"be/neurade/v2/internal/repository"
	"be/neurade/v2/internal/util"
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB             *gorm.DB
	UserRepository *repository.UserRepository
	Log            *logrus.Logger
}

func NewUserService(db *gorm.DB, userRepository *repository.UserRepository, log *logrus.Logger) *UserService {
	return &UserService{
		DB:             db,
		UserRepository: userRepository,
		Log:            log,
	}
}

func (s *UserService) Update(ctx context.Context, request *model.UserUpdateRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx)
	defer tx.Rollback()

	user := &entity.User{
		ID:          request.ID,
		Email:       request.Email,
		Role:        request.Role,
		GithubToken: request.GithubToken,
	}
	if err := s.UserRepository.Update(tx, user); err != nil {
		s.Log.Error("Failed to update user: ", err)
		return nil, err
	}
	return converter.UserToResponse(user), nil
}

func (s *UserService) Register(ctx context.Context, request *model.UserCreateRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	existingUser, _ := s.UserRepository.FindUserByEmail(request.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Error("Failed to hash password: ", err)
		return nil, err
	}

	user := &entity.User{
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		Role:         request.Role,
		Verified:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.UserRepository.Create(tx, user); err != nil {
		tx.Rollback()
		s.Log.Error("Failed to create user: ", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Error("Failed to commit transaction: ", err)
		return nil, err
	}
	createdUser, err := s.UserRepository.FindUserByEmail(user.Email)
	if err != nil {
		tx.Rollback()
		s.Log.Error("Failed to find created user: ", err)
		return nil, err
	}

	return converter.UserToResponse(createdUser), nil
}

func (s *UserService) Login(ctx context.Context, request *model.LoginRequest, jwtUtil *util.JWTUtil) (*model.AuthResponse, error) {
	user, err := s.UserRepository.FindUserByEmail(request.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		s.Log.Error("Failed to generate token: ", err)
		return nil, err
	}

	return &model.AuthResponse{
		Token: token,
		User:  *converter.UserToResponse(user),
	}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx)
	users := &[]entity.User{}

	err := s.UserRepository.FindAll(tx, users)
	if err != nil {
		s.Log.Error("Failed to get all users: ", err)
		return nil, err
	}

	var userResponses []*model.UserResponse
	for _, user := range *users {
		userResponses = append(userResponses, converter.UserToResponse(&user))
	}

	return userResponses, nil
}
