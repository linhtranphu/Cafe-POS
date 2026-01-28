package user

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleWaiter   Role = "waiter"
	RoleCashier  Role = "cashier"
	RoleManager  Role = "manager"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"-"`
	Role      Role               `bson:"role" json:"role"`
	Name      string             `bson:"name" json:"name"`
	Active    bool               `bson:"active" json:"active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	LastLogin *time.Time         `bson:"last_login,omitempty" json:"last_login,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}