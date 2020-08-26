package bprometheus

import (
	"net/http"

	"github.com/go-masonry/mortar/providers/groups"
	"github.com/go-masonry/mortar/providers/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

// PrometheusInternalHandlerFxOption fx.Provide option that will register Prometheus
// HTTP handler to serve "/metrics" endpoint on internal port
func PrometheusInternalHandlerFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  groups.InternalHTTPHandlers,
		Target: PrometheusHTTPHandlerPatternPair,
	})
}

// PrometheusHTTPHandlerPatternPair provides mortar Internal HTTP Pattern Pair
// It can later be registered to serve metrics endpoint on internal port
func PrometheusHTTPHandlerPatternPair() types.HTTPHandlerPatternPair {
	return types.HTTPHandlerPatternPair{
		Pattern: "/metrics",
		Handler: HTTPHandler(),
	}
}

// HTTPHandler provides the Prometheus HTTP scrape handler.
func HTTPHandler() http.Handler {
	return promhttp.Handler()
}
