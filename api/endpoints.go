package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/AheadAviation/bagshop-catalog/item"
)

type Endpoints struct {
	CreateItemEndpoint endpoint.Endpoint
	GetItemsEndpoint   endpoint.Endpoint
	HealthEndpoint     endpoint.Endpoint
}

func MakeEndpoints(s Service, tracer stdopentracing.Tracer) Endpoints {
	return Endpoints{
		CreateItemEndpoint: opentracing.TraceServer(tracer, "POST /api/v1/catalog/items")(MakeCreateItemEndpoint(s)),
		GetItemsEndpoint:   opentracing.TraceServer(tracer, "GET /api/v1/catalog/items")(MakeGetItemsEndpoint(s)),
		HealthEndpoint:     opentracing.TraceServer(tracer, "GET /api/v1/catalog/health")(MakeHealthEndpoint(s)),
	}
}

func MakeCreateItemEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var span stdopentracing.Span
		span, ctx = stdopentracing.StartSpanFromContext(ctx, "create item")
		span.SetTag("service", "catalog")
		defer span.Finish()
		req := request.(createItemRequest)
		id, err := s.CreateItem(req.Name, req.Description, req.Price, req.Count)
		return postResponse{ID: id}, err
	}
}

func MakeGetItemsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var span stdopentracing.Span
		span, ctx = stdopentracing.StartSpanFromContext(ctx, "get items")
		span.SetTag("service", "catalog")
		defer span.Finish()
		its, err := s.GetItems()
		return itemsResponse{Items: its}, err
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

type createItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Count       int     `json:"count"`
}

type postResponse struct {
	ID string `json:"id"`
}

type itemsResponse struct {
	Items []item.Item `json:"items"`
}

type healthRequest struct {
	// pass
}

type healthResponse struct {
	Health []Health `json:"health"`
}
