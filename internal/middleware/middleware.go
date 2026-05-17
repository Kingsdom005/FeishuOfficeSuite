package middleware

import (
	"context"
	"log"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func Logging(logger log.Logger) transport.ServerMiddleware {
	return func(next transport.Handler) transport.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				startTime = time.Now()
				kind      = "server"
				operation = ""
			)

			if tr, ok := transport.FromServerContext(ctx); ok {
				operation = tr.Operation()
				kind = tr.Kind().String()
			}

			logHelper := log.NewHelper(logger)
			logHelper.WithContext(ctx).Infof("[%s] %s | %s | started", kind, operation, time.Since(startTime))

			reply, err = next(ctx, req)

			logHelper.WithContext(ctx).Infof("[%s] %s | %s | %s | error: %v",
				kind, operation, time.Since(startTime), statusCode(err), err)

			return
		}
	}
}

func Tracing(tracer trace.Tracer) transport.ServerMiddleware {
	return func(next transport.Handler) transport.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var span trace.Span
			ctx, span = tracer.Start(ctx, "server",
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()

			if tr, ok := transport.FromServerContext(ctx); ok {
				span.SetAttributes(
					attribute.String("rpc.system", "grpc"),
					attribute.String("rpc.service", tr.Service()),
					attribute.String("rpc.method", tr.Operation()),
				)
			}

			propagator := otel.GetTextMapPropagator()
			metadataClient := propagator.Extract(ctx, propagation.HeaderCarrier{})

			span.SetAttributes(
				attribute.String("request.id", getRequestID(ctx)),
			)

			reply, err = next(ctx, req)

			if err != nil {
				span.RecordError(err)
			}

			return reply, err
		}
	}
}

func Metrics() transport.ServerMiddleware {
	return func(next transport.Handler) transport.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			startTime := time.Now()

			reply, err = next(ctx, req)

			duration := time.Since(startTime).Milliseconds()

			var operation string
			if tr, ok := transport.FromServerContext(ctx); ok {
				operation = tr.Operation()
			}

			log.Printf("metrics: operation=%s duration_ms=%d status=%s",
				operation, duration, statusCode(err))

			return reply, err
		}
	}
}

func CORS() transport.ServerMiddleware {
	return func(next transport.Handler) transport.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			return next(ctx, req)
		}
	}
}

func statusCode(err error) string {
	if err != nil {
		return "error"
	}
	return "ok"
}

func getRequestID(ctx context.Context) string {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return ""
	}
	if reqID := md.Get("x-request-id"); len(reqID) > 0 {
		return reqID[0]
	}
	return ""
}

type tracing struct {
	tracer trace.Tracer
}