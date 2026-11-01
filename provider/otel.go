//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package provider

import (
	"context"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type otelProvider struct {
	traceProvider *tracesdk.TracerProvider
}

func (p *otelProvider) Shutdown(ctx context.Context) error {
	return p.traceProvider.Shutdown(ctx)
}
