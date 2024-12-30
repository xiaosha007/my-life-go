package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	model "github.com/xiaosha007/my-life-go/internal/models"
	"github.com/xiaosha007/my-life-go/pkg/crypto"
)

type LoginResponse struct {
	User  *LoginUserResponse
	Token string
}

type LoginUserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type IAuthService interface {
	Login(username, password string) (*LoginResponse, error)
	// Logout() bool
	CreateToken(username string) (string, error)
	VerifyToken(tokenString string) error
}

type IAuthUserService interface {
	GetByUsername(username string) (*model.User, error)
}

type AuthService struct {
	userService IAuthUserService
	jwtSecret   string
}

func NewAuthService(userService IAuthUserService, jwtSecret string) *AuthService {
	return &AuthService{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

func (a *AuthService) Login(username, password string) (*LoginResponse, error) {
	// get user from db by username
	user, err := a.userService.GetByUsername(username)

	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errors.New("error fetching user: user not found")
		}
		return nil, err
	}

	// compare the password
	if !crypto.IsSame(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// generate jwt token if valid
	token, err := a.CreateToken(username)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User: &LoginUserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		},
		Token: token,
	}, nil
}

func (a *AuthService) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AuthService) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).Before(time.Now()) {
		return errors.New("token has expired")
	}

	return nil
}
