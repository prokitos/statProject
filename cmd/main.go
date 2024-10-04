package main

import (
	"myMod/internal/app"
	"myMod/internal/clickhouseStat"
	"myMod/internal/config"
	"myMod/internal/mysql"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("log is loaded")

	var cfg config.MainConfig
	cfg.ConfigMustLoad("docker")
	log.Debug("config is loaded")

	clickDB, err := mysql.NewClickhouseWriter("127.0.0.1", 9000, "rotator", "statistics", "default", "qwerty123")
	if err != nil {
		return
	}

	time.Sleep(time.Second)

	statsManager := clickhouseStat.NewManager(clickDB, 10*time.Second)
	statsManager.StartTimer()

	var application app.App
	application.NewManager(statsManager)
	application.NewServer(cfg.Server.Port)
	log.Debug("server is loaded")

}
