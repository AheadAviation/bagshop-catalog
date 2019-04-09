package main

import (
	"fmt"
	corelog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	_ "github.com/go-sql-driver/mysql"
	flag "github.com/namsral/flag"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"github.com/AheadAviation/bagshop-catalog/api"
	"github.com/AheadAviation/bagshop-catalog/db"
	"github.com/AheadAviation/bagshop-catalog/db/mysql"
)

const ServiceName = "catalog"

var (
	port string
)

var (
	HTTPLatency = stdprometheus.NewHistogramVec(stdprometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Time (in seconds) spent serving HTTP requests.",
		Buckets: stdprometheus.DefBuckets,
	}, []string{"method", "path", "status_code", "isWS"})
)

func init() {
	stdprometheus.MustRegister(HTTPLatency)
	flag.StringVar(&port, "catalog-port", "8083", "Port on which to run the catalog service")
	db.Register("mysql", &mysql.MySQL{})
}

func main() {
	flag.Parse()
	errc := make(chan error)

	// Logging Domain
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Tracing Domain
	var tracer stdopentracing.Tracer
	{
		tracer = stdopentracing.NoopTracer{}
		stdopentracing.InitGlobalTracer(tracer)
	}

	// Database Domain
	dbconn := false
	for !dbconn {
		err := db.Init()
		if err != nil {
			if err == db.ErrNoDatabaseSelected {
				corelog.Fatal(err)
			}
			corelog.Print(err)
		} else {
			dbconn = true
		}
	}

	fieldKeys := []string{"method"}

	// Service Domain
	var service api.Service
	{
		service = api.NewFixedService()
		service = api.LoggingMiddleware(logger)(service)
		service = api.NewInstrumentingService(
			kitprometheus.NewCounterFrom(
				stdprometheus.CounterOpts{
					Namespace: "microservices_demo",
					Subsystem: "catalog",
					Name:      "request_count",
					Help:      "Number of requests received",
				},
				fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "microservices_demo",
				Subsystem: "catalog",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			service,
		)
	}

	// Endpoint Domain
	endpoints := api.MakeEndpoints(service, tracer)

	// Transport Domain
	router := api.MakeHTTPHandler(endpoints, logger, tracer)

	// HTTP Server
	go func() {
		logger.Log("transport", "HTTP", "port", port)
		errc <- http.ListenAndServe(fmt.Sprintf(":%v", port), router)
	}()

	// Capture Interrupts
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("exit", <-errc)
}
