//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

// Package trace provides utility functions and types
// for easy instrumenting of applications.
package trace

import (
	"context"
	"reflect"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracePoint wraps a [trace.Tracer] object to support
// tracing of an applications control point.
type TracePoint struct {
	name   string
	tracer trace.Tracer
}

// NewTracePoint create a new [TracePoint] using the given type argument as
// the tracer name, the given name argument as the default span name and the
// given options as is.
//
// In case the default span name is empty, the span name is automatically
// derived from the call hierarchy.
func NewTracePoint[T any](name string, opts ...trace.TracerOption) *TracePoint {
	return &TracePoint{
		name:   name,
		tracer: otel.Tracer(reflect.TypeFor[T]().PkgPath(), opts...),
	}
}

// StartInternal invokes [trace.Tracer.Start] with [trace.SpanKindInternal]
// as well as the given attributes.
func (t *TracePoint) StartInternal(ctx context.Context, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.start(ctx, trace.WithSpanKind(trace.SpanKindInternal), trace.WithAttributes(attributes...))
}

// StartServer invokes [trace.Tracer.Start] with [trace.SpanKindServer]
// as well as the given attributes.
func (t *TracePoint) StartServer(ctx context.Context, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.start(ctx, trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attributes...))
}

// StartClient invokes [trace.Tracer.Start] with [trace.SpanKindClient]
// as well as the given attributes.
func (t *TracePoint) StartClient(ctx context.Context, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.start(ctx, trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attributes...))
}

// StartProducer invokes [trace.Tracer.Start] with [trace.SpanKindProducer]
// as well as the given attributes.
func (t *TracePoint) StartProducer(ctx context.Context, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.start(ctx, trace.WithSpanKind(trace.SpanKindProducer), trace.WithAttributes(attributes...))
}

// StartConsumer invokes [trace.Tracer.Start] with [trace.SpanKindConsumer]
// as well as the given attributes.
func (t *TracePoint) StartConsumer(ctx context.Context, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.start(ctx, trace.WithSpanKind(trace.SpanKindConsumer), trace.WithAttributes(attributes...))
}

func (t *TracePoint) start(ctx context.Context, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	name := t.name
	if name == "" {
		name = callerName()
	}
	return t.tracer.Start(ctx, name, opts...)
}

func callerName() string {
	pc, _, _, _ := runtime.Caller(3)
	name := runtime.FuncForPC(pc).Name()
	split := strings.Split(name, "/")
	return split[len(split)-1]
}

// SetStatus sets the given span's status based on the given error.
func SetStatus(span trace.Span, err error) error {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	return err
}
