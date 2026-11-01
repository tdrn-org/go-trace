//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package trace_test

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tdrn-org/go-config-toml"
	"github.com/tdrn-org/go-trace"
	"github.com/tdrn-org/go-trace/provider"
)

func TestTracing(t *testing.T) {
	provider := startTestProvider(t)
	defer provider.Shutdown(t.Context())

	tracer := trace.NewTracePoint[testing.T]("")
	traceCtx, span := tracer.StartClient(t.Context())
	defer span.End()

	nestedFunction(traceCtx, tracer, 10)
}

func nestedFunction(ctx context.Context, tracer *trace.TracePoint, depth int) {
	if depth > 0 {
		traceCtx, span := tracer.StartInternal(ctx)
		defer span.End()

		nestedFunction(traceCtx, tracer, depth-1)
	}
}

func startTestProvider(t *testing.T) provider.Provider {
	protocol := os.Getenv("TRACING_PROTOCOL")
	endpoint := os.Getenv("TRACING_ENDPOINT_URL")
	if protocol == "" || endpoint == "" {
		t.Skip("tracing protocol/endpoint url not set; skipping tests")
	}
	endpointURL, err := url.Parse(endpoint)
	require.NoError(t, err)
	cfg := &provider.TracingConfig{
		Enabled:       true,
		Protocol:      provider.TracingProtocol(protocol),
		EndpointURL:   config.URLSpec{URL: endpointURL},
		BatchTimeout:  config.DurationSpec(5 * time.Second),
		ExportTimeout: config.DurationSpec(30 * time.Second),
	}
	provider, err := provider.Start(t.Context(), cfg)
	require.NoError(t, err)
	return provider
}
