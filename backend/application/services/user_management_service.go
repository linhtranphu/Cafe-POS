package services

import (
	"context"
	"errors"
	"time"
	"cafe-pos/backend/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserManagementService struct {
	userRepo    UserRepository
	authService *AuthService
}

func NewUserManagementService(userRepo UserRepository, authService *AuthService) *UserManagementService {
	return &UserManagementService{
		userRepo:    userRepo,
		authService: authService,
	}
}

type CreateUserRequest struct {
	Username string    `json:"username" binding:"required,min=3,max=50"`
	Password string    `json:"password" binding:"required,min=6"`
	Name     string    `json:"name" binding:"required,min=2,max=100"`
	Role     user.Role `json:"role" binding:"required"`
	Active   bool      `json:"active"`
}

type UpdateUserRequest struct {
	Name   string    `json:"name" binding:"required,min=2,max=100"`
	Role   user.Role `json:"role" binding:"required"`
	Active bool      `json:"active"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Role      user.Role `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}

func (s *UserManagementService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	// Check if username already exists
	existingUser, _ := s.userRepo.FindByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Validate role
	if !isValidRole(req.Role) {
		return nil, errors.New("invalid role")
	}

	// Hash password
	hashedPassword, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	newUser := &user.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Name:      req.Name,
		Role:      req.Role,
		Active:    req.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return s.toUserResponse(newUser), nil
}

func (s *UserManagementService) GetAllUsers(ctx context.Context) ([]*UserResponse, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []*UserResponse
	for _, u := range users {
		responses = append(responses, s.toUserResponse(u))
	}

	return responses, nil
}

func (s *UserManagementService) GetUser(ctx context.Context, id string) (*UserResponse, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.toUserResponse(u), nil
}

func (s *UserManagementService) UpdateUser(ctx context.Context, id string, req *UpdateUserRequest) (*UserResponse, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate role
	if !isValidRole(req.Role) {
		return nil, errors.New("invalid role")
	}

	// Update user fields
	u.Name = req.Name
	u.Role = req.Role
	u.Active = req.Active
	u.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, userID, u); err != nil {
		return nil, err
	}

	return s.toUserResponse(u), nil
}

func (s *UserManagementService) ResetPassword(ctx context.Context, id string, req *ResetPasswordRequest) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Hash new password
	hashedPassword, err := s.authService.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	u.Password = hashedPassword
	u.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, userID, u)
}

func (s *UserManagementService) ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	u, err := s.userRepo.FindByID(ctx, userObjID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify current password
	if !s.authService.CheckPassword(req.CurrentPassword, u.Password) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := s.authService.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	u.Password = hashedPassword
	u.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, userObjID, u)
}

func (s *UserManagementService) ToggleUserStatus(ctx context.Context, id string) (*UserResponse, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	u.Active = !u.Active
	u.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, userID, u); err != nil {
		return nil, err
	}

	return s.toUserResponse(u), nil
}

func (s *UserManagementService) DeleteUser(ctx context.Context, id string) error {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Prevent deleting the last manager
	if u.Role == user.RoleManager {
		managers, _ := s.userRepo.FindByRole(ctx, user.RoleManager)
		if len(managers) <= 1 {
			return errors.New("cannot delete the last manager")
		}
	}

	return s.userRepo.Delete(ctx, userID)
}

func (s *UserManagementService) GetUsersByRole(ctx context.Context, role user.Role) ([]*UserResponse, error) {
	users, err := s.userRepo.FindByRole(ctx, role)
	if err != nil {
		return nil, err
	}

	var responses []*UserResponse
	for _, u := range users {
		responses = append(responses, s.toUserResponse(u))
	}

	return responses, nil
}

func (s *UserManagementService) GetActiveUsers(ctx context.Context) ([]*UserResponse, error) {
	users, err := s.userRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	var responses []*UserResponse
	for _, u := range users {
		responses = append(responses, s.toUserResponse(u))
	}

	return responses, nil
}

func (s *UserManagementService) toUserResponse(u *user.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID.Hex(),
		Username:  u.Username,
		Name:      u.Name,
		Role:      u.Role,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		LastLogin: u.LastLogin,
	}
}

func isValidRole(role user.Role) bool {
	validRoles := []user.Role{user.RoleManager, user.RoleCashier, user.RoleWaiter}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}