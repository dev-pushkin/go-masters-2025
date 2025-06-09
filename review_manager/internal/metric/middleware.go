package metric

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/prometheus/client_golang/prometheus"
)

// return http middlware for get metric and store it
func StoreMetricsMiddleware(basePath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			metric.reqInFlightGauge.Inc()
			defer metric.reqInFlightGauge.Dec()

			m := httpsnoop.CaptureMetrics(next, w, r)
			path := strings.TrimPrefix(r.URL.Path, basePath)
			httpStatus := strconv.Itoa(m.Code)

			// count total request
			l := requestCountLabel(path, r.Method, httpStatus)
			metric.reqTotalCounterVec.With(l).Inc()

			// histogram for request duration
			l = requestMetricLabel(path, r.Method, httpStatus)
			metric.reqDurationSecHistogram.With(l).Observe(m.Duration.Seconds())
		})
	}
}

func requestMetricLabel(path, method, httpStatus string) prometheus.Labels {
	return prometheus.Labels{
		routePathLabel:  path,
		httpMethodLabel: method,
		httpStatusLabel: httpStatus,
	}
}

func requestCountLabel(path, method, httpStatus string) prometheus.Labels {
	return prometheus.Labels{
		routePathLabel:  path,
		httpMethodLabel: method,
		httpStatusLabel: httpStatus,
	}
}
