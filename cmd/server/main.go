package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/feishu/feishu-office-suite/internal/data"
	"github.com/feishu/feishu-office-suite/internal/job"
	"github.com/feishu/feishu-office-suite/internal/handler"
	"github.com/feishu/feishu-office-suite/internal/middleware"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf ./configs")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server, rr registry.Registrar, dp registry.Discovery) *kratos.App {
	return kratos.New(
		kratos.Name("feishu-office-suite"),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
		kratos.Registrar(rr),
		kratos.Discoverer(dp),
	)
}

func main() {
	flag.Parse()

	logger := log.NewStdLogger(os.Stdout)
	log := log.NewHelper(logger)

	cfg, err := config.NewConfigWithPath(flagconf)
	if err != nil {
		log.Errorf("Failed to load config: %v", err)
		os.Exit(1)
	}

	bc := &data.BootstrapConfig{}
	if err := cfg.Get("bootstrap").Scan(bc); err != nil {
		log.Errorf("Failed to scan bootstrap config: %v", err)
	}

	otelTracerProvider, err := initTracer(bc.OpenTelemetry)
	if err != nil {
		log.Warnf("Failed to init tracer: %v", err)
	}

	tracer := otel.Tracer("feishu-office-suite")

	dataData, cleanupData, err := data.NewData(log, bc, otelTracerProvider)
	if err != nil {
		log.Errorf("Failed to create data: %v", err)
		os.Exit(1)
	}
	defer cleanupData()

	jobExecutor := job.NewExecutor(dataData)
	asynqWorker := job.NewAsynqWorker(jobExecutor)

	go func() {
		if err := asynqWorker.Start(); err != nil {
			log.Errorf("Failed to start asynq worker: %v", err)
		}
	}()

	userHandler := handler.NewUserHandler(dataData)
	messageHandler := handler.NewMessageHandler(dataData)
	calendarHandler := handler.NewCalendarHandler(dataData)
	approvalHandler := handler.NewApprovalHandler(dataData)

	httpServer := http.NewServer(
		http.Middleware(
			recovery.Recovery(),
			middleware.Logging(log),
			middleware.Metrics(),
			middleware.Tracing(tracer),
			middleware.CORS(),
		),
	)
	httpServer.HandlePrefix("/", userHandler)
	httpServer.HandlePrefix("/", messageHandler)
	httpServer.HandlePrefix("/", calendarHandler)
	httpServer.HandlePrefix("/", approvalHandler)

	grpcServer := grpc.NewServer(
		grpc.Middleware(
			recovery.Recovery(),
			middleware.Logging(log),
			middleware.Tracing(tracer),
		),
	)
	handler.RegisterFeishuUserServer(grpcServer, userHandler)
	handler.RegisterFeishuMessageServer(grpcServer, messageHandler)
	handler.RegisterFeishuCalendarServer(grpcServer, calendarHandler)
	handler.RegisterFeishuApprovalServer(grpcServer, approvalHandler)

	app, cleanup, err := wireApp(
		log,
		cfg,
		httpServer,
		grpcServer,
		dataData,
	)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func initTracer(cfg data.OpenTelemetryConfig) (*sdktrace.TracerProvider, error) {
	if cfg.Endpoint == "" {
		return nil, nil
	}

	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("feishu-office-suite"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}