package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	model "github.com/xiaosha007/my-life-go/internal/models"
	"github.com/xiaosha007/my-life-go/internal/services"
)

const minPasswordLength, maxPasswordLength = 5, 30

type UserHandler struct {
	Service services.IUserService
}

func NewUserHandler(service services.IUserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	user, err := h.Service.GetByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passwordLength := len(user.Password)
	if passwordLength < minPasswordLength || passwordLength > maxPasswordLength {
		http.Error(w, "Invalid password length!", http.StatusBadRequest)
		return
	}

	err = h.Service.Create(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)

}
