package main

import (
	"github.com/xiaosha007/my-life-go/internal/config"
	"github.com/xiaosha007/my-life-go/internal/db"
)

func main() {
	// init config
	config.Init()
	configData := config.GetConfig()

	// init db
	defer db.CloseConnection()

	db.InitConnection(configData.DB)

}
