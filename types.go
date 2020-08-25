package bprometheus

import (
	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/prometheus/client_golang/prometheus"
)

type promCounterVec struct {
	*prometheus.CounterVec
}

type promCounter struct {
	prometheus.Counter
}

type promGaugeVec struct {
	*prometheus.GaugeVec
}

type promGauge struct {
	prometheus.Gauge
}

type promHistogramVec struct {
	*prometheus.HistogramVec
}

type promHistogram struct {
	prometheus.Observer
}

func (p *promCounterVec) WithTags(tags map[string]string) monitor.Counter {
	return &promCounter{
		Counter: p.With(tags),
	}
}

func (p *promGaugeVec) WithTags(tags map[string]string) monitor.Gauge {
	return &promGauge{
		Gauge: p.With(tags),
	}
}

func (p *promHistogramVec) WithTags(tags map[string]string) monitor.Histogram {
	return &promHistogram{
		Observer: p.With(tags),
	}
}

func (p *promHistogram) Record(v float64) {
	p.Observe(v)
}

var _ monitor.BricksCounter = (*promCounterVec)(nil)
var _ monitor.Counter = (*promCounter)(nil)
var _ monitor.BricksGauge = (*promGaugeVec)(nil)
var _ monitor.Gauge = (*promGauge)(nil)
var _ monitor.BricksHistogram = (*promHistogramVec)(nil)
var _ monitor.Histogram = (*promHistogram)(nil)
