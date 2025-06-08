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
			semconv.ServiceNameKey.String("service_a"),
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

	tracer = otel.Tracer("service_a")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "service_a_handler")
		defer span.End()

		logger.Info().Str("trace_id", span.SpanContext().TraceID().String()).Msg("Handling request in Service A")

		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8081", nil)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create request to Service B")
			http.Error(w, "Failed to create request to Service B", http.StatusInternalServerError)
			return
		}

		client := otelhttp.DefaultClient // to link traces with service b
		resp, err := client.Do(req)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to call Service B")
			http.Error(w, "Failed to call Service B", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to read response from Service B")
			http.Error(w, "Failed to read response from Service B", http.StatusInternalServerError)
			return
		}

		io.WriteString(w, "Crnogorski rap is rolling on Service A!\n")
		io.WriteString(w, string(body))
	})

	logger.Info().Msg("Service A is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
} 