package service

import (
	"context"
	"errors"
	"note-api/auth"
	"note-api/entity"
	"note-api/repository"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort   = errors.New("password must be at least 8 charecters")
	ErrEmailRequired      = errors.New("email is required")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password string) (*entity.User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, ErrEmailRequired
	}
	if len(password) < 8 {
		return nil, ErrPasswordTooShort
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}
	user := &entity.User{
		Email:        strings.TrimSpace(email),
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}
	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil

}

type LoginResult struct {
	User  *entity.User
	Token string
}

func (s *UserService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, ErrEmailRequired
	}
	if len(password) < 8 {
		return nil, ErrPasswordTooShort
	}
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		User:  user,
		Token: token,
	}, nil

}
