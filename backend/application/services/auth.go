package services

import (
	"context"
	"errors"
	"cafe-pos/backend/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*user.User, error)
	FindAll(ctx context.Context) ([]*user.User, error)
	FindByRole(ctx context.Context, role user.Role) ([]*user.User, error)
	FindActive(ctx context.Context) ([]*user.User, error)
	Create(ctx context.Context, user *user.User) error
	Update(ctx context.Context, id primitive.ObjectID, user *user.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	UpdateLastLogin(ctx context.Context, id primitive.ObjectID) error
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

	// Update last login
	a.userRepo.UpdateLastLogin(ctx, u.ID)

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

func (a *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}