// Package metrics Prometheus 指标收集
package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

var (
	// HTTPRequestsTotal HTTP 请求总数
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gohub_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestDuration HTTP 请求耗时
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gohub_http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"method", "path"},
	)

	// CacheHitsTotal 缓存命中总数
	CacheHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gohub_cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	// CacheMissesTotal 缓存未命中总数
	CacheMissesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gohub_cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)

	// DBQueryDuration 数据库查询耗时
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gohub_db_query_duration_seconds",
			Help:    "Database query latencies in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"operation"},
	)

	// ActiveConnections 活跃连接数
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "gohub_active_connections",
			Help: "Number of active connections",
		},
	)

	// ErrorsTotal 错误总数
	ErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gohub_errors_total",
			Help: "Total number of errors",
		},
		[]string{"error_type"},
	)
)

// PrometheusMiddleware Prometheus 指标收集中间件
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// 增加活跃连接
		ActiveConnections.Inc()
		defer ActiveConnections.Dec()

		// 处理请求
		c.Next()

		// 记录指标
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		duration := time.Since(start).Seconds()

		HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)

		// 如果有错误，记录错误指标
		if len(c.Errors) > 0 {
			ErrorsTotal.WithLabelValues("http_error").Inc()
		}
	}
}

// Handler 返回 Prometheus 指标处理器
func Handler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// RecordCacheHit 记录缓存命中
func RecordCacheHit(cacheType string) {
	CacheHitsTotal.WithLabelValues(cacheType).Inc()
}

// RecordCacheMiss 记录缓存未命中
func RecordCacheMiss(cacheType string) {
	CacheMissesTotal.WithLabelValues(cacheType).Inc()
}

// RecordDBQuery 记录数据库查询耗时
func RecordDBQuery(operation string, duration time.Duration) {
	DBQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordError 记录错误
func RecordError(errorType string) {
	ErrorsTotal.WithLabelValues(errorType).Inc()
}
