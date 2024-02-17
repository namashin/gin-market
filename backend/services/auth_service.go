package services

import (
	"fmt"
	"gin-market/models"
	"gin-market/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type IAuthService interface {
	SignUp(email string, password string) error
	Login(email string, password string) (*string, error)
	Logout(token string) error
	GetUserFromToken(tokenString string) (*models.User, error)
	InvalidateToken(tokenString string) error
	IsValidToken(tokenString string) bool
}

type AuthService struct {
	repository repository.IAuthRepository
}

func NewAuthService(repository repository.IAuthRepository) IAuthService {
	return &AuthService{
		repository: repository,
	}
}

func (as *AuthService) SignUp(email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// make user
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	return as.repository.CreateUser(user)
}

func (as *AuthService) Login(email string, password string) (*string, error) {
	targetUser, err := as.repository.FindUser(email)
	if err != nil {
		return nil, err
	}

	// ? check password
	err = bcrypt.CompareHashAndPassword([]byte(targetUser.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	token, err := CreateToken(targetUser.ID, targetUser.Email)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (as *AuthService) Logout(token string) error {
	return as.repository.InvalidateToken(token)
}

func CreateToken(userId uint, email string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("secret_key")))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (as *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("secret_key")), nil
	})
	if err != nil {
		return nil, err
	}

	var user *models.User

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, jwt.ErrHashUnavailable
		}

		user, err = as.repository.FindUser(claims["email"].(string))
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (as *AuthService) InvalidateToken(tokenString string) error {
	return as.repository.InvalidateToken(tokenString)
}

func (as *AuthService) IsValidToken(tokenString string) bool {
	return as.repository.IsValidToken(tokenString)
}
