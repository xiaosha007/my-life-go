package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gorilla/mux"
	model "github.com/xiaosha007/my-life-go/internal/models"
	"github.com/xiaosha007/my-life-go/internal/services"
	"github.com/xiaosha007/my-life-go/pkg/utils"
)

const minPasswordLength, maxPasswordLength = 5, 30

type UserHandler struct {
	Service      services.IUserService
	StatsdClient *statsd.Client
}

func NewUserHandler(service services.IUserService, statsdClient *statsd.Client) *UserHandler {
	return &UserHandler{
		Service:      service,
		StatsdClient: statsdClient,
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	h.StatsdClient.Count("get_user_by_id.count", 1, []string{}, 1)

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	user, err := h.Service.GetByID(id)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, *user)

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
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid password length!")
		return
	}

	err = h.Service.Create(&user)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, user)

}
