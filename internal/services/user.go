package services

import (
	"log"

	"github.com/xiaosha007/my-life-go/internal/db"
	model "github.com/xiaosha007/my-life-go/internal/models"
	"github.com/xiaosha007/my-life-go/pkg/crypto"
)

type IUserService interface {
	GetByID(id int) (*model.User, error)
	Create(user *model.User) error
	GetByUsername(username string) (*model.User, error)
}

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetByID(id int) (*model.User, error) {
	var user model.User

	err := db.DB.Get(&user, "SELECT * FROM users WHERE id = $1", id)

	return &user, err
}

func (s *UserService) GetByUsername(username string) (*model.User, error) {
	var user model.User

	err := db.DB.Get(&user, "SELECT * FROM users WHERE username = $1", username)

	return &user, err
}

func (s *UserService) Create(user *model.User) error {

	encryptedPassword, err := crypto.Encrypt(user.Password)

	if err != nil {
		log.Println("failed to encrypt password")
		return err
	}

	user.Password = encryptedPassword

	query := "INSERT INTO users (name, username, password, email) VALUES ($1 ,$2, $3, $4) RETURNING id"

	return db.DB.QueryRow(query, user.Name, user.Username, user.Password, user.Email).Scan(&user.ID)
}
