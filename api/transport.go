package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sony/gobreaker"
)

func MakeHTTPHandler(e Endpoints, logger log.Logger, tracer stdopentracing.Tracer) *mux.Router {
	r := mux.NewRouter().StrictSlash(false)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// GET /api/v1/catalog/health Health Check
	// GET /api/v1/catalog/metrics Prometheus-style metrics

	r.Methods("GET").PathPrefix("/api/v1/catalog/health").Handler(httptransport.NewServer(
		circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Health",
			Timeout: 30 * time.Second,
		}))(e.HealthEndpoint),
		decodeHealthRequest,
		encodeHealthResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GET /api/v1/catalog/health", logger)))...,
	))
	r.Handle("/api/v1/catalog/metrics", promhttp.Handler())
	return r
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	switch err {
	case ErrNotFound:
		code = http.StatusNotFound
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":       err.Error(),
		"status_code": code,
		"status_text": http.StatusText(code),
	})
}

func decodeHealthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return struct{}{}, nil
}

func encodeHealthResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return encodeResponse(ctx, w, response.(healthResponse))
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
