package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/xiaosha007/my-life-go/internal/config"
	"github.com/xiaosha007/my-life-go/internal/db"
	"github.com/xiaosha007/my-life-go/internal/monitoring"
	"github.com/xiaosha007/my-life-go/internal/router"
)

func main() {
	// init config
	config.Init()
	configData := config.GetConfig()

	// init db
	defer db.CloseConnection()

	db.InitConnection(configData.DB)

	statsdClient := monitoring.NewStatsdClient(configData.Statsd)

	r := router.NewRouter(statsdClient, configData.JwtSecret)

	port := 8080

	log.Printf("Starting server on port %s", strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":8080", r))

}
