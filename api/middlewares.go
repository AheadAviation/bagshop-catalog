package api

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"

	"github.com/AheadAviation/bagshop-catalog/item"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) CreateItem(name, description string, price float32,
	count int) (id string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "CreateItem",
			"item_name", name,
			"item_id", id,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.CreateItem(name, description, price, count)
}

func (mw loggingMiddleware) GetItems() (its []item.Item, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetItems",
			"result", len(its),
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.GetItems()
}

func (mw loggingMiddleware) Health() (health []Health) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Health",
			"result", len(health),
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Health()
}

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(requestCount metrics.Counter, requestLatency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   requestCount,
		requestLatency: requestLatency,
		Service:        s,
	}
}

func (s *instrumentingService) CreateItem(name, description string, price float32,
	count int) (id string, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "create_item").Add(1)
		s.requestLatency.With("method", "create_item").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.CreateItem(name, description, price, count)
}

func (s *instrumentingService) GetItems() ([]item.Item, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get_items").Add(1)
		s.requestLatency.With("method", "get_items").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetItems()
}

func (s *instrumentingService) Health() []Health {
	defer func(begin time.Time) {
		s.requestCount.With("method", "health").Add(1)
		s.requestLatency.With("method", "health").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Health()
}
