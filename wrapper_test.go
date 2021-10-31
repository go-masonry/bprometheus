package bprometheus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/suite"
)

type promSuite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(promSuite))
}

func (d *promSuite) TestNamespace() {
	namespace := "service"
	reporter := Builder().SetNamespace(namespace).Build()
	err := reporter.Connect(context.Background())
	d.NoError(err)
	counter, err := reporter.Metrics().Counter("counter_namespace", "")
	d.NoError(err)
	counterWithTags, err := counter.WithTags(nil)
	d.NoError(err)
	counterWithTags.Inc()
	err = reporter.Close(context.Background())
	d.NoError(err)
	handler := HTTPHandler()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, &http.Request{})
	body := recorder.Body.String()
	d.Require().Contains(body, "service_counter")
}

func (d *promSuite) TestTags() {
	tagKeys := []string{"one", "three"}
	tags := map[string]string{
		"one":   "two",
		"three": "four",
	}
	reporter := Builder().Build()
	err := reporter.Connect(context.Background())
	d.NoError(err)
	counter, err := reporter.Metrics().Counter("counter_tags", "", tagKeys...)
	d.NoError(err)
	counterWithTags, err := counter.WithTags(tags)
	d.NoError(err)
	counterWithTags.Inc()
	err = reporter.Close(context.Background())
	d.NoError(err)
	handler := HTTPHandler()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, &http.Request{})
	body := recorder.Body.String()
	d.Require().Contains(body, "one=\"two\"")
	d.Require().Contains(body, "three=\"four\"")
}

func (d *promSuite) TestCounter() {
	reporter := Builder().Build()
	err := reporter.Connect(context.Background())
	d.NoError(err)
	counter, err := reporter.Metrics().Counter("counter", "")
	d.NoError(err)
	counterWithTags, err := counter.WithTags(nil)
	d.NoError(err)
	counterWithTags.Add(5)
	err = reporter.Close(context.Background())
	d.NoError(err)
	handler := HTTPHandler()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, &http.Request{})
	body := recorder.Body.String()
	d.Require().Contains(body, "counter 5")
}

func (d *promSuite) TestGauge() {
	reporter := Builder().Build()
	err := reporter.Connect(context.Background())
	d.NoError(err)
	gauge, err := reporter.Metrics().Gauge("gauge", "")
	d.NoError(err)
	gaugeWithTags, err := gauge.WithTags(nil)
	d.NoError(err)
	gaugeWithTags.Add(5)
	err = reporter.Close(context.Background())
	d.NoError(err)
	handler := HTTPHandler()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, &http.Request{})
	body := recorder.Body.String()
	d.Require().Contains(body, "gauge 5")
}

func (d *promSuite) TestHistogram() {
	reporter := Builder().Build()
	err := reporter.Connect(context.Background())
	d.NoError(err)
	histogram, err := reporter.Metrics().Histogram("histogram", "", nil)
	d.NoError(err)
	histogramWithTags, err := histogram.WithTags(nil)
	d.NoError(err)
	histogramWithTags.Record(7)
	err = reporter.Close(context.Background())
	d.NoError(err)
	handler := HTTPHandler()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, &http.Request{})
	body := recorder.Body.String()
	d.Require().Contains(body, "histogram_sum 7")
	d.Require().Contains(body, "histogram_count 1")
}

func (d *promSuite) TestTimer() {
	reporter := Builder().Build()
	err := reporter.Connect(context.Background())
	d.NoError(err)
	timer, err := reporter.Metrics().Timer("timer", "")
	d.NoError(err)
	timerWithTags, err := timer.WithTags(nil)
	d.NoError(err)
	timerWithTags.Record(7 * time.Second)
	timerWithTags.Record(100 * time.Second)
	timerWithTags.Record(1000 * time.Second)
	err = reporter.Close(context.Background())
	d.NoError(err)
	handler := HTTPHandler()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, &http.Request{})
	body := recorder.Body.String()
	d.Require().Contains(body, "timer_sum 1107")
	d.Require().Contains(body, "timer_count 3")
}

func (d *promSuite) TestRemove() {
	reporter := Builder().Build()
	err := reporter.Connect(context.Background())
	d.NoError(err)
	counter, err := reporter.Metrics().Counter("counter_remove", "")
	d.NoError(err)
	err = reporter.Metrics().Remove(counter)
	d.NoError(err)
}

func (d *promSuite) TestCustomCollectors() {
	counterCollector := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "custom",
		Name:      "counter",
		Help:      "",
	})
	anotherCollector := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "custom",
		Name:      "another",
		Help:      "",
	})
	reporter := Builder().AddPredefinedCollectors(counterCollector).AddPredefinedCollectors(anotherCollector).SetNamespace("different").Build()
	err := reporter.Connect(context.Background())
	d.NoError(err)
	// Count
	counterCollector.Inc()
	anotherCollector.Inc()

	handler := HTTPHandler()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, &http.Request{})
	body := recorder.Body.String()
	d.Require().Contains(body, "counter 1")
	d.Require().Contains(body, "another 1")
}
