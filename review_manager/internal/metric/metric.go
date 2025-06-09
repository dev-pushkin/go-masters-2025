package metric

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// Metric contains general metrics for api methods.
var metric struct {
	reqInFlightGauge        prometheus.Gauge
	reqTotalCounterVec      *prometheus.CounterVec
	reqDurationSecHistogram *prometheus.HistogramVec
}

const (
	routePathLabel  = "path"
	httpMethodLabel = "method"
	httpStatusLabel = "status"
)

var (
	httpRequestDurationBuckets = []float64{0.1, 0.5, 1, 2.5, 5, 10}
)

func InitAndServe(opt ...Options) error {
	o := newDefaultOptions()
	for _, option := range opt {
		option(o)
	}

	if o.regestry == nil {
		return fmt.Errorf("prometheus registry is not set")
	}

	InitMetrics(o.regestry, o.namespace, o.subsystem)

	addr := fmt.Sprintf("%s:%d", o.host, o.port)
	go func() {
		if err := ServeMetrics(o.regestry, addr, o.path); err != nil {
			panic(fmt.Sprintf("failed to start metrics server: %v", err))
		}
	}()
	return nil
}

func InitMetrics(reg *prometheus.Registry, namespace, subsystem string) {
	metric.reqInFlightGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_requests_in_flight",
			Help:      "Amount of currently processing API requests at moment.",
		},
	)
	reg.MustRegister(metric.reqInFlightGauge)

	metric.reqTotalCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_requests_total",
			Help:      "Amount of processed API requests.",
		},
		[]string{routePathLabel, httpMethodLabel, httpStatusLabel},
	)
	reg.MustRegister(metric.reqTotalCounterVec)

	metric.reqDurationSecHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_request_duration_seconds",
			Help:      "API request latency distributions.",
			Buckets:   httpRequestDurationBuckets,
		},
		[]string{routePathLabel, httpMethodLabel, httpStatusLabel},
	)
	reg.MustRegister(metric.reqDurationSecHistogram)

	// Collector with GO metrics
	goCollector := collectors.NewGoCollector()
	reg.MustRegister(goCollector)

	// Collector with process metrics
	processCollector := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{Namespace: namespace})
	reg.MustRegister(processCollector)
}
