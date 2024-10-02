package transport

import (
	"myMod/internal/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// здесь хранятся хэндлеры.

func SetHandlers(instance *gin.Engine) {

	reg := prometheus.NewRegistry()

	reg.MustRegister(metrics.RequestDuration)
	reg.MustRegister(metrics.RequestStatus)

	instance.Use(metrics.Observer)

	instance.POST("/task", insertTask)
	instance.GET("/task/:id", getTask)
	instance.DELETE("/task/:id", deleteTask)
	instance.PUT("/task", updateTask)

	instance.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
