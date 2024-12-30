package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	model "github.com/xiaosha007/my-life-go/internal/models"
	"github.com/xiaosha007/my-life-go/internal/services"
	"github.com/xiaosha007/my-life-go/pkg/crypto"
)

// Mock for AuthUserService
type mockAuthUserService struct {
	mock.Mock
}

func (m *mockAuthUserService) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)

	// Check if args.Get(0) is nil before type assertion
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.User), args.Error(1)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockUserService := new(mockAuthUserService)
	authService := services.NewAuthService(mockUserService, "test-secret")

	password, _ := crypto.Encrypt("password123")

	// Mock user data
	mockUser := &model.User{
		ID:       "123",
		Name:     "Test User",
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: password,
	}
	mockUserService.On("GetByUsername", "testuser").Return(mockUser, nil)

	// Call Login
	response, err := authService.Login("testuser", "password123")

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "123", response.User.ID)
	assert.Equal(t, "Test User", response.User.Name)
	assert.NotEmpty(t, response.Token)
	mockUserService.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockUserService := new(mockAuthUserService)
	authService := services.NewAuthService(mockUserService, "test-secret")

	password, _ := crypto.Encrypt("password123")

	// Mock user data
	mockUser := &model.User{
		ID:       "123",
		Name:     "Test User",
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: password,
	}
	mockUserService.On("GetByUsername", "testuser").Return(mockUser, nil)

	// Call Login with invalid password
	response, err := authService.Login("testuser", "wrongpassword")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.EqualError(t, err, "invalid credentials")
	mockUserService.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockUserService := new(mockAuthUserService)
	authService := services.NewAuthService(mockUserService, "test-secret")

	// Mock GetByUsername to return an error
	mockUserService.On("GetByUsername", "testuser").Return(nil, errors.New("user not found"))

	// Call Login
	response, err := authService.Login("testuser", "password123")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.EqualError(t, err, "user not found")
	mockUserService.AssertExpectations(t)
}

func TestAuthService_CreateToken(t *testing.T) {
	authService := services.NewAuthService(nil, "test-secret")

	// Call CreateToken
	token, err := authService.CreateToken("testuser")

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse the token to validate it
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, "testuser", claims["username"])
	assert.NotZero(t, claims["exp"])
}

func TestAuthService_VerifyToken_Valid(t *testing.T) {
	authService := services.NewAuthService(nil, "test-secret")

	// Generate a token
	token, _ := authService.CreateToken("testuser")

	// Verify the token
	err := authService.VerifyToken(token)

	// Assertions
	assert.NoError(t, err)
}

func TestAuthService_VerifyToken_Invalid(t *testing.T) {
	authService := services.NewAuthService(nil, "test-secret")

	// Invalid token string
	invalidToken := "invalid.token.string"

	// Verify the token
	err := authService.VerifyToken(invalidToken)

	// Assertions
	assert.Error(t, err)
}

func TestAuthService_VerifyToken_Expired(t *testing.T) {
	authService := services.NewAuthService(nil, "test-secret")

	// Create an expired token
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, _ := expiredToken.SignedString([]byte("test-secret"))

	// Verify the token
	err := authService.VerifyToken(tokenString)

	// Assertions
	assert.Error(t, err)
}
