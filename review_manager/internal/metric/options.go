package metric

import "github.com/prometheus/client_golang/prometheus"

type Options func(*o)

type o struct {
	regestry  *prometheus.Registry
	host      string
	port      int
	path      string // path for metrics endpoint
	namespace string
	subsystem string
}

func newDefaultOptions() *o {
	return &o{
		regestry:  prometheus.NewPedanticRegistry(),
		host:      "localhost",
		port:      9090,
		path:      "/metrics",
		namespace: "app",
		subsystem: "api",
	}
}

func WithRegistry(reg *prometheus.Registry) Options {
	return func(o *o) {
		o.regestry = reg
	}
}
func WithHost(host string) Options {
	return func(o *o) {
		o.host = host
	}
}
func WithPort(port int) Options {
	return func(o *o) {
		o.port = port
	}
}
func WithPath(path string) Options {
	return func(o *o) {
		o.path = path
	}
}
func WithNamespace(namespace string) Options {
	return func(o *o) {
		o.namespace = namespace
	}
}
func WithSubsystem(subsystem string) Options {
	return func(o *o) {
		o.subsystem = subsystem
	}
}
