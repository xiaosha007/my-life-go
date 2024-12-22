package router

import (
	"github.com/DataDog/datadog-go/statsd"
	"github.com/gorilla/mux"
	v1 "github.com/xiaosha007/my-life-go/api/v1"
	"github.com/xiaosha007/my-life-go/internal/services"
)

func NewRouter(statsdClient *statsd.Client) *mux.Router {
	r := mux.NewRouter()

	userService := services.UserService{}
	userHandler := v1.NewUserHandler(&userService, statsdClient)

	// API routes
	r.HandleFunc("/api/v1/users/{id}", userHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/api/v1/users", userHandler.CreateUser).Methods("POST")

	return r
}
