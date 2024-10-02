package transport

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func metricGet(c *gin.Context) {

}

var requestMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "clean",
	Subsystem:  "http",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"status"})

var methodMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "clean",
	Subsystem:  "http",
	Name:       "methods",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"method"})

var routeMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "clean",
	Subsystem:  "http",
	Name:       "routes",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"route"})

func ObserveRequest(durat time.Duration, status int, method string, name string) {
	requestMetric.WithLabelValues(strconv.Itoa(status)).Observe(durat.Seconds())
	methodMetric.WithLabelValues(method).Observe(1)
	routeMetric.WithLabelValues(name).Observe(1)
}
