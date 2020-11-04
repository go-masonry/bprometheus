package bprometheus

import (
	"context"
	"errors"

	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	errInvalidMetricType = errors.New("invalid metric type")
)

type promWrapper struct {
	namespace  string
	collectors []prometheus.Collector
}

func newPromWrapper(cfg *promConfig) monitor.BricksReporter {
	return &promWrapper{
		namespace:  cfg.namespace,
		collectors: cfg.collectors,
	}
}

func (p *promWrapper) Connect(ctx context.Context) error {
	// Register all predefined Collectors
	for _, collector := range p.collectors {
		if err := prometheus.Register(collector); err != nil {
			return err
		}
	}
	return nil
}

func (p *promWrapper) Close(ctx context.Context) error {
	return nil
}

func (p *promWrapper) Metrics() monitor.BricksMetrics {
	return p
}

func (p *promWrapper) Counter(name, desc string, tagKeys ...string) (monitor.BricksCounter, error) {
	counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: p.namespace,
		Name:      name,
		Help:      desc,
	}, tagKeys)
	err := prometheus.Register(counterVec)
	if err != nil {
		return nil, err
	}
	return &promCounterVec{
		CounterVec: counterVec,
	}, nil
}

func (p *promWrapper) Gauge(name, desc string, tagKeys ...string) (monitor.BricksGauge, error) {
	gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: p.namespace,
		Name:      name,
		Help:      desc,
	}, tagKeys)
	err := prometheus.Register(gaugeVec)
	if err != nil {
		return nil, err
	}
	return &promGaugeVec{
		GaugeVec: gaugeVec,
	}, nil
}

func (p *promWrapper) Histogram(name, desc string, buckets []float64, tagKeys ...string) (monitor.BricksHistogram, error) {
	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: p.namespace,
		Name:      name,
		Help:      desc,
		Buckets:   buckets,
	}, tagKeys)
	err := prometheus.Register(histogramVec)
	if err != nil {
		return nil, err
	}
	return &promHistogramVec{
		HistogramVec: histogramVec,
	}, nil
}

func (p *promWrapper) Timer(name, desc string, tagKeys ...string) (monitor.BricksTimer, error) {
	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: p.namespace,
		Name:      name,
		Help:      desc,
	}, tagKeys)
	err := prometheus.Register(histogramVec)
	if err != nil {
		return nil, err
	}
	return &promTimerVec{
		HistogramVec: histogramVec,
	}, nil
}

func (p *promWrapper) Remove(metric monitor.BrickMetric) error {
	collector, ok := metric.(prometheus.Collector)
	if !ok {
		return errInvalidMetricType
	}
	prometheus.Unregister(collector)
	return nil
}

var _ monitor.BricksReporter = (*promWrapper)(nil)
