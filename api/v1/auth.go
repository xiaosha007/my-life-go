package v1

import (
	"encoding/json"
	"net/http"

	model "github.com/xiaosha007/my-life-go/internal/models"
	"github.com/xiaosha007/my-life-go/internal/services"
	"github.com/xiaosha007/my-life-go/pkg/utils"
)

type AuthHandler struct {
	Service services.IAuthService
}

func NewAuthHandler(authService services.IAuthService) *AuthHandler {
	return &AuthHandler{
		Service: authService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user *model.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	loginResponse, err := h.Service.Login(user.Username, user.Password)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, loginResponse)
}
