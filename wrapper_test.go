package bprometheus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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
	counter.WithTags(nil).Inc()
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
	counter.WithTags(tags).Inc()
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
	counter.WithTags(nil).Add(5)
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
	gauge.WithTags(nil).Add(5)
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
	histogram.WithTags(nil).Record(7)
	err = reporter.Close(context.Background())
	d.NoError(err)
	handler := HTTPHandler()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, &http.Request{})
	body := recorder.Body.String()
	d.Require().Contains(body, "histogram_sum 7")
	d.Require().Contains(body, "histogram_count 1")
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
