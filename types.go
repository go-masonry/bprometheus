package bprometheus

import (
	"time"

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

type promTimerVec struct {
	*prometheus.HistogramVec
}

type promTimer struct {
	prometheus.Observer
}

func (p *promCounterVec) WithTags(tags map[string]string) (monitor.Counter, error) {
	counter, err := p.GetMetricWith(tags)
	if err != nil {
		return nil, err
	}
	return &promCounter{
		Counter: counter,
	}, nil
}

func (p *promGaugeVec) WithTags(tags map[string]string) (monitor.Gauge, error) {
	gauge, err := p.GetMetricWith(tags)
	if err != nil {
		return nil, err
	}
	return &promGauge{
		Gauge: gauge,
	}, nil
}

func (p *promHistogramVec) WithTags(tags map[string]string) (monitor.Histogram, error) {
	histogram, err := p.GetMetricWith(tags)
	if err != nil {
		return nil, err
	}
	return &promHistogram{
		Observer: histogram,
	}, nil
}

func (p *promHistogram) Record(v float64) {
	p.Observe(v)
}

func (p *promTimerVec) WithTags(tags map[string]string) (monitor.Timer, error) {
	histogram, err := p.GetMetricWith(tags)
	if err != nil {
		return nil, err
	}
	return &promTimer{
		Observer: histogram,
	}, nil
}

func (p *promTimer) Record(d time.Duration) {
	p.Observe(d.Seconds())
}

var _ monitor.BricksCounter = (*promCounterVec)(nil)
var _ monitor.Counter = (*promCounter)(nil)
var _ monitor.BricksGauge = (*promGaugeVec)(nil)
var _ monitor.Gauge = (*promGauge)(nil)
var _ monitor.BricksHistogram = (*promHistogramVec)(nil)
var _ monitor.Histogram = (*promHistogram)(nil)
var _ monitor.Timer = (*promTimer)(nil)
