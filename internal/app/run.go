package app

import (
	"fmt"
	"myMod/internal/database"

	"github.com/gin-gonic/gin"
)

// запуск и остановка сервера.

type App struct {
	Server *gin.Engine
	Stats  *database.Manager
}

func (a *App) NewManager(manager *database.Manager) {
	a.Stats = manager
}

func (a *App) NewServer(port string) {
	app := gin.New()
	a.Server = app
	a.SetHandlers()
	a.launchServer(port)
}

func (a *App) Stop() {
	fmt.Println("Gracefully shutting down...")
	//a.Server.(50 * time.Second)
}

func (a *App) launchServer(port string) {
	a.Server.Run(port)
}
