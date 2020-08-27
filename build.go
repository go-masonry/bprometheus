package bprometheus

import (
	"container/list"

	"github.com/go-masonry/mortar/interfaces/monitor"
)

type promConfig struct {
	namespace string
}
type promBuilder struct {
	ll *list.List
}

// PrometheusBuilder defines Prometheus builder
type PrometheusBuilder interface {
	monitor.Builder
	SetNamespace(namespace string) monitor.Builder
}

// Builder creates a builder to create Prometheus client
func Builder() PrometheusBuilder {
	return &promBuilder{
		ll: list.New(),
	}
}

func (s *promBuilder) SetNamespace(namespace string) monitor.Builder {
	s.ll.PushBack(func(cfg *promConfig) {
		cfg.namespace = namespace
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
