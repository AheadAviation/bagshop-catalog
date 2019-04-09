package api

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
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

func (s *instrumentingService) Health() []Health {
	defer func(begin time.Time) {
		s.requestCount.With("method", "health").Add(1)
		s.requestLatency.With("method", "health").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Health()
}
