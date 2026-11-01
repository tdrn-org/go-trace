//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

// Package provider contains convienence functions and types
// to support easy setup of tracing.
package provider

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

// Provider interface is used to access the configured tracing provider.
type Provider interface {
	// Shutdown shuts down the tracing provider.
	Shutdown(ctx context.Context) error
}

// Start setups and starts tracing using the given configuration.
func Start(ctx context.Context, cfg *TracingConfig) (Provider, error) {
	if !cfg.Enabled {
		return &noopProvider{}, nil
	}
	exporter, err := newExporter(ctx, cfg)
	if err != nil {
		return nil, err
	}
	serviceName := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(cfg.ServiceName))
	traceProvider := tracesdk.NewTracerProvider(
		withBatcher(exporter, time.Duration(cfg.BatchTimeout), time.Duration(cfg.ExportTimeout)),
		tracesdk.WithResource(serviceName))
	otel.SetTracerProvider(traceProvider)
	provider := &otelProvider{
		exporter:      exporter,
		traceProvider: traceProvider,
	}
	slog.Info("tracing enabled")
	return provider, nil
}

func withBatcher(exporter *otlptrace.Exporter, batchTimeout time.Duration, exportTimeout time.Duration) tracesdk.TracerProviderOption {
	opts := []tracesdk.BatchSpanProcessorOption{
		tracesdk.WithBatchTimeout(batchTimeout),
		tracesdk.WithExportTimeout(exportTimeout),
	}
	return tracesdk.WithBatcher(exporter, opts...)
}

func newExporter(ctx context.Context, cfg *TracingConfig) (*otlptrace.Exporter, error) {
	switch cfg.Protocol {
	case TracingProtocolHttp:
		return newHttpExporter(ctx, cfg.EndpointURL.URL, cfg.Headers)
	case TracingProtocolGRPC:
		return newGRPCExporter(ctx, cfg.EndpointURL.URL, cfg.Headers)
	default:
		return nil, fmt.Errorf("unrecognized tracing protocol: %q", cfg.Protocol)
	}
}

func newHttpExporter(ctx context.Context, endpoint *url.URL, headers map[string]string) (*otlptrace.Exporter, error) {
	opts := make([]otlptracehttp.Option, 0)
	opts = append(opts, otlptracehttp.WithEndpointURL(endpoint.String()))
	if endpoint.Scheme == "http" {
		slog.Warn("using insecure tracing endpoint URL", slog.Any("endpoint", endpoint))
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	opts = append(opts, otlptracehttp.WithHeaders(headers))
	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create http Exporter (cause: %w)", err)
	}
	return exporter, nil
}

func newGRPCExporter(ctx context.Context, endpoint *url.URL, headers map[string]string) (*otlptrace.Exporter, error) {
	opts := make([]otlptracegrpc.Option, 0)
	opts = append(opts, otlptracegrpc.WithEndpointURL(endpoint.String()))
	if endpoint.Scheme == "http" {
		slog.Warn("using insecure tracing gRPC endpoint", slog.Any("endpoint", endpoint))
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	opts = append(opts, otlptracegrpc.WithHeaders(headers))
	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC Exporter (cause: %w)", err)
	}
	return exporter, nil
}
