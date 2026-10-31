//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package trace

import (
	"context"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/trace"
)

type TracePoint struct {
	tracer trace.Tracer
}

func NewTracePoint[T any](opts ...trace.TracerOption) *TracePoint {
	return &TracePoint{
		tracer: otel.Tracer(reflect.TypeFor[T]().PkgPath(), opts...),
	}
}

func (t *TracePoint) StartInternal(ctx context.Context, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.start(ctx, trace.SpanKindInternal, attributes...)
}

func (t *TracePoint) StartServer(ctx context.Context, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.start(ctx, trace.SpanKindServer, attributes...)
}

func (p *TracePoint) StartHttpServer(req *http.Request, attributes ...attribute.KeyValue) (*http.Request, trace.Span) {
	traceCtx, span := p.start(req.Context(), trace.SpanKindServer, attributes...)
	SetHttpRequest(span, req)
	return req.WithContext(traceCtx), span
}

func (t *TracePoint) StartClient(ctx context.Context, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.start(ctx, trace.SpanKindClient, attributes...)
}

func (t *TracePoint) StartProducer(ctx context.Context, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.start(ctx, trace.SpanKindProducer, attributes...)
}

func (t *TracePoint) StartConsuming(ctx context.Context, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.start(ctx, trace.SpanKindConsumer, attributes...)
}

func (t *TracePoint) start(ctx context.Context, kind trace.SpanKind, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, callerName(), trace.WithSpanKind(kind), trace.WithAttributes(attributes...))
}

func callerName() string {
	pc, _, _, _ := runtime.Caller(3)
	name := runtime.FuncForPC(pc).Name()
	split := strings.Split(name, "/")
	return split[len(split)-1]
}

func SetHttpRequest(span trace.Span, req *http.Request) {
	span.SetAttributes(attribute.String(string(semconv.HTTPRequestMethodKey), req.Method))
	span.SetAttributes(attribute.String(string(semconv.URLPathKey), req.URL.Path))
	span.SetAttributes(attribute.String(string(semconv.URLSchemeKey), req.URL.Scheme))
}

func SetHttpResponse(span trace.Span, rsp *http.Response) {
	span.SetAttributes(attribute.Int(string(semconv.HTTPResponseStatusCodeKey), rsp.StatusCode))
}

func SetStatus(span trace.Span, err error) error {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	return err
}
