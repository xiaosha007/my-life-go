package router

import (
	"github.com/DataDog/datadog-go/statsd"
	"github.com/gorilla/mux"
	v1 "github.com/xiaosha007/my-life-go/api/v1"
	"github.com/xiaosha007/my-life-go/internal/services"
)

func NewRouter(statsdClient *statsd.Client, jwtSecret string) *mux.Router {
	r := mux.NewRouter()

	userService := services.NewUserService()
	userHandler := v1.NewUserHandler(userService, statsdClient)

	authService := services.NewAuthService(userService, jwtSecret)
	authHandler := v1.NewAuthHandler(authService)

	// API routes
	r.HandleFunc("/api/v1/users/{id}", userHandler.GetUserByID).Methods("GET")
	r.HandleFunc("/api/v1/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")

	return r
}
