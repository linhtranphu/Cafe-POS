package services

import (
	"errors"
	"time"
	"cafe-pos/backend/domain/user"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey string
}

type Claims struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Role     user.Role `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey string) *JWTService {
	return &JWTService{secretKey: secretKey}
}

func (j *JWTService) GenerateToken(u *user.User) (string, error) {
	claims := Claims{
		UserID:   u.ID.Hex(),
		Username: u.Username,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}