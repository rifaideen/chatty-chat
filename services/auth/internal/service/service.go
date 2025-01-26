package service

import (
	"errors"
	"pkg/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	Login(request *auth.LoginRequest) (string, error)
	VerifyToken(request *auth.VerifyRequest) (*jwt.Token, error)
}

type AuthService struct {
	expiry int    // expiry in days
	secret string // secret key
}

func New(secret string, expiry int) Service {
	return &AuthService{
		expiry: expiry,
		secret: secret,
	}
}

func (s *AuthService) Login(request *auth.LoginRequest) (string, error) {
	if request.Username == "admin" && request.Password == "admin" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  1,
			"username": "admin",
			"exp":      time.Now().Add(time.Hour * 24 * time.Duration(s.expiry)).Unix(),
		})

		return token.SignedString([]byte(s.secret))
	}

	return "", errors.New("invalid credentials")
}

func (s *AuthService) VerifyToken(request *auth.VerifyRequest) (*jwt.Token, error) {
	return jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
}
