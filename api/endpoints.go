package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	stdopentracing "github.com/opentracing/opentracing-go"
)

type Endpoints struct {
	HealthEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service, tracer stdopentracing.Tracer) Endpoints {
	return Endpoints{
		HealthEndpoint: opentracing.TraceServer(tracer, "GET /health")(MakeHealthEndpoint(s)),
	}
}

func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var span stdopentracing.Span
		span, ctx = stdopentracing.StartSpanFromContext(ctx, "health check")
		span.SetTag("service", "catalog")
		defer span.Finish()
		health := s.Health()
		return healthResponse{Health: health}, nil
	}
}

type healthRequest struct {
	// pass
}

type healthResponse struct {
	Health []Health `json:"health"`
}
