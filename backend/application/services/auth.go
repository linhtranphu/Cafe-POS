package services

import (
	"context"
	"errors"
	"cafe-pos/backend/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	Create(ctx context.Context, user *user.User) error
}

type AuthService struct {
	userRepo   UserRepository
	jwtService *JWTService
}

func NewAuthService(userRepo UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (a *AuthService) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	u, err := a.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !u.Active {
		return nil, errors.New("user is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := a.jwtService.GenerateToken(u)
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		Token: token,
		User:  *u,
	}, nil
}

func (a *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}