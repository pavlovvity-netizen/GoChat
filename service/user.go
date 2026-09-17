package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"GoChat/domain"
)

type UserRepository interface {
	RegisterUser(user *domain.User) error
	GetUserByName(name string) (*domain.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(user *domain.RegisterRequest) (*domain.User, error) {
	passwordHash, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		Username:     user.Username,
		PasswordHash: passwordHash,
	}

	err = s.repo.RegisterUser(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) GetUserByName(req *domain.LoginRequest) (*domain.User, error) {
	user, err := s.repo.GetUserByName(req.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("неверное имя пользователя или пароль")
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
