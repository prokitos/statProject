package main

import (
	"myMod/internal/app"
	"myMod/internal/config"
	"myMod/internal/database"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("log is loaded")

	var cfg config.MainConfig
	cfg.ConfigMustLoad("docker")
	log.Debug("config is loaded")

	var CDB database.ClickDatabase
	CDB.ClickHouseStart()

	manager := database.NewManager(&CDB, 10*time.Second)
	manager.StartTimer()

	var application app.App
	application.NewManager(manager)
	application.NewServer(cfg.Server.Port)
	log.Debug("server is loaded")

}
