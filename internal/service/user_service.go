package service

import (
	"errors"
	"user-management/internal/domain"
	"user-management/internal/repository"
	"user-management/pkg/security"
)

type UserService interface {
	Register(user *domain.User) error
	Login(email, password string) (*domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) Register(user *domain.User) error {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.FindByUsername(user.Username)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	existingUser, _ = s.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// 密码加密
	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepo.Create(user)
}

func (s *userService) Login(email, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := security.CheckPasswordHash(password, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// 其他方法实现...
