package app

import (
	"myMod/internal/database"
	"myMod/internal/metrics"
	"myMod/internal/transport"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// здесь хранятся хэндлеры.

func (a *App) SetHandlers() {
	instance := a.Server

	reg := prometheus.NewRegistry()

	reg.MustRegister(metrics.RequestDuration)
	reg.MustRegister(metrics.RequestStatus)

	instance.Use(metrics.Observer)
	instance.Use(a.GetClickMetrics)

	instance.POST("/task", transport.InsertTask)
	instance.GET("/task/:id", transport.GetTask)
	instance.DELETE("/task/:id", transport.DeleteTask)
	instance.PUT("/task", transport.UpdateTask)

	instance.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))
}

func (a *App) GetClickMetrics(ctx *gin.Context) {

	ua := string(ctx.Request.UserAgent())
	parsed := user_agent.New(ua)
	browserName, _ := parsed.Browser()

	var data database.Statistic
	data.Os = parsed.OS()
	data.Browser = browserName
	data.Request = 1
	data.Country = "russia"
	currentTime := time.Now()
	data.Timestamp = currentTime.Format("2006.01.02 15:04:05")

	defer func() {
		a.Stats.Adding(data)
	}()

	if time.Now().Second() > 30 {
		data.Impression++
	}

}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
