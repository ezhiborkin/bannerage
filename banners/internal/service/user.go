package service

import (
	"banners/domain/models"
	"banners/lib/logger/sl"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserStorage interface {
	CreateUserStorage(email, role string, passHash []byte) error
	GetUserStorage(email string) (*models.User, error)
}

func (s *Service) CreateUser(email, role string, password string) error {
	const op = "service.CreateUser"

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to generate passHash", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.userStorage.CreateUserStorage(email, role, passHash)
	if err != nil {
		s.log.Error("failed to create a user", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) LoginUser(email, password string) (string, error) {
	const op = "service.LoginUser"

	expiresAt := time.Now().Add(time.Minute * 100).Unix()

	user, err := s.userStorage.GetUserStorage(email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	tk := &models.Token{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}
