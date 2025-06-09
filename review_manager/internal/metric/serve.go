package metric

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ServeMetrics(reg *prometheus.Registry, addr, path string) error {
	http.Handle(path, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	return http.ListenAndServe(addr, nil)
}
