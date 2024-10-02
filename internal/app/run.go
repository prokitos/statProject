package app

import (
	"fmt"
	"myMod/internal/transport"

	"github.com/gin-gonic/gin"
)

// запуск и остановка сервера.

type App struct {
	Server *gin.Engine
}

func (a *App) NewServer(port string) {
	app := gin.New()
	a.Server = app
	a.setHandler()
	a.launchServer(port)
}

func (a *App) Stop() {
	fmt.Println("Gracefully shutting down...")
	//a.Server.(50 * time.Second)
}

func (a *App) setHandler() {
	transport.SetHandlers(a.Server)
}

func (a *App) launchServer(port string) {
	a.Server.Run(port)
}
