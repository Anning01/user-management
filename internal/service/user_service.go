package service

import (
	"errors"

	"github.com/Anning01/user-management/internal/domain"
	"github.com/Anning01/user-management/internal/repository"
	"github.com/Anning01/user-management/pkg/security"
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

func (s *userService) GetUserByID(id uint) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) UpdateUser(user *domain.User) error {
	// 如果更新了邮箱，检查是否已被其他用户使用
	existingUser, _ := s.userRepo.FindByEmail(user.Email)
	if existingUser != nil && existingUser.ID != user.ID {
		return errors.New("email already exists")
	}

	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}
