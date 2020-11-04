package bprometheus

import (
	"container/list"

	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/prometheus/client_golang/prometheus"
)

type promConfig struct {
	namespace  string
	collectors []prometheus.Collector
}
type promBuilder struct {
	ll *list.List
}

// PrometheusBuilder defines Prometheus builder
type PrometheusBuilder interface {
	monitor.Builder
	// SetNamespace allows to set a default Prometheus Namespace
	SetNamespace(namespace string) PrometheusBuilder
	// AddCollectors allows to register predefined Collectors to the same Prometheus Registry
	// *** Actual registration will occur only when the `monitor.BricksReporter.Connect(ctx)` is called ***
	// *** Any error returned during registration will fail the `Connect` method.
	AddCollectors(collectors ...prometheus.Collector) PrometheusBuilder
}

// Builder creates a builder to create Prometheus client
func Builder() PrometheusBuilder {
	return &promBuilder{
		ll: list.New(),
	}
}

func (s *promBuilder) SetNamespace(namespace string) PrometheusBuilder {
	s.ll.PushBack(func(cfg *promConfig) {
		cfg.namespace = namespace
	})
	return s
}

func (s *promBuilder) AddCollectors(collectors ...prometheus.Collector) PrometheusBuilder {
	s.ll.PushBack(func(cfg *promConfig) {
		cfg.collectors = append(cfg.collectors, collectors...)
	})
	return s
}

func (s *promBuilder) Build() monitor.BricksReporter {
	cfg := &promConfig{}
	if s != nil {
		for e := s.ll.Front(); e != nil; e = e.Next() {
			f := e.Value.(func(config *promConfig))
			f(cfg)
		}
	}
	return newPromWrapper(cfg)
}

var _ monitor.Builder = (*promBuilder)(nil)
