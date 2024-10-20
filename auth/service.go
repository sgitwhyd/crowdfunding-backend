package auth

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
)

type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type service struct {}

func NewService() *service {
	return &service{}
}

var SECRET_KEY = []byte(os.Getenv("JWT_SECRET"))

func (s *service) GenerateToken(userID int) (string, error) {

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})

	signedToken, err := t.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (s *service) ValidateToken(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return t, err
	}

	return t, nil
}