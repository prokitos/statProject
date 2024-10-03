package app

import (
	"myMod/internal/clickhouseStat"
	"myMod/internal/metrics"
	"myMod/internal/transport"

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

	statsKey := clickhouseStat.NewKey(clickhouseStat.Key{
		Os:      parsed.OS(),
		Browser: browserName,
	})

	statsValue := clickhouseStat.Value{Requests: 1}

	defer func() {
		a.Stats.Append(statsKey, statsValue)
	}()

	// country := "russia"
	// statsKey.Country = country

	// if time.Now().Second() > 30 {
	// 	statsValue.Impressions++
	// }

}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
