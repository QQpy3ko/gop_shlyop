package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var tracer trace.Tracer

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("service_b"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp, nil
}

func main() {
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	tp, err := initTracer()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize tracer")
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Fatal().Err(err).Msg("Failed to shutdown tracer provider")
		}
	}()

	tracer = otel.Tracer("service_b")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), "service_b_handler")
		defer span.End()

		logger = logger.With().Str("trace_id", span.SpanContext().TraceID().String()).Logger()

		for name, values := range r.Header {
			for _, value := range values {
				logger.Info().Str("header", name).Str("value", value).Msg("Incoming header")
			}
		}

		logger.Info().Msg("Handling request in Service B")
		io.WriteString(w, "Crnogorski rap is rolling on Service B!\n")
	})

	http.Handle("/", otelhttp.NewHandler(handler, "service_b_http_handler")) // to link traces with service a

	logger.Info().Msg("Service B is running on :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}