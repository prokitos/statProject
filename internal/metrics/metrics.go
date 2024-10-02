package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// var requestMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
// 	Namespace:  "clean",
// 	Subsystem:  "http",
// 	Name:       "request",
// 	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
// }, []string{"status"})

// var methodMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
// 	Namespace:  "clean",
// 	Subsystem:  "http",
// 	Name:       "methods",
// 	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
// }, []string{"method"})

// var routeMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
// 	Namespace:  "clean",
// 	Subsystem:  "http",
// 	Name:       "routes",
// 	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
// }, []string{"route"})

// func ObserveRequest(durat time.Duration, status int, method string, name string) {
// 	requestMetric.WithLabelValues(strconv.Itoa(status)).Observe(durat.Seconds())
// 	methodMetric.WithLabelValues(method).Observe(1)
// 	routeMetric.WithLabelValues(name).Observe(1)
// }

// Define metrics
// requestCount := prometheus.NewCounterVec(
// 	prometheus.CounterOpts{
// 		Name: "http_requests_total",
// 		Help: "Total number of HTTP requests",
// 	},
// 	[]string{"method", "path"},
// )
// var requestMetric = prometheus.NewSummaryVec(prometheus.SummaryOpts{
// 	Namespace:  "clean",
// 	Subsystem:  "http",
// 	Name:       "request",
// 	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
// }, []string{"status"})

var RequestDuration *prometheus.SummaryVec = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "http_request_duration_seconds",
		Help:       "Duration of HTTP requests in seconds",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)
var RequestStatus *prometheus.SummaryVec = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "http_status_total",
		Help:       "Total number of HTTP requests and status",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"status"},
)

func Observer(c *gin.Context) {
	start := time.Now()
	c.Next()
	duration := time.Since(start).Seconds()

	RequestStatus.WithLabelValues(strconv.Itoa(c.Writer.Status())).Observe(1)
	RequestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(duration)
}
