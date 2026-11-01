//
// Copyright (C) 2026 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package trace

import (
	"net/http"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/trace"
)

// StartHttpServer invokes [trace.Tracer.Start] with [trace.SpanKindServer]
// as well as the given attributes. The necessary context is derived from the
// request argument which in turn is wrapped with the created trace context.
//
// The created span is already populated with the HTTP attributes from the
// given request (see [SetHttpRequest]).
func (p *TracePoint) StartHttpServer(req *http.Request, attributes ...attribute.KeyValue) (*http.Request, trace.Span) {
	traceCtx, span := p.start(req.Context(), trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attributes...), trace.WithAttributes(httpRequestAttributes(req)...))
	return req.WithContext(traceCtx), span
}

func httpRequestAttributes(req *http.Request) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.HTTPRequestMethodKey.String(req.Method),
		semconv.URLPath(req.URL.Path),
		semconv.URLScheme(req.URL.Scheme),
	}
}

// SetHttpRequest populates the given span with the HTTP attributes
// from the given request.
func SetHttpRequest(span trace.Span, req *http.Request) {
	if !span.IsRecording() {
		return
	}

	for _, attribute := range httpRequestAttributes(req) {
		span.SetAttributes(attribute)
	}
}

// SetHttpResponse populates the given span with the HTTP attributes
// from the given response.
func SetHttpResponse(span trace.Span, rsp *http.Response) {
	if !span.IsRecording() {
		return
	}

	span.SetAttributes(semconv.HTTPResponseStatusCode(rsp.StatusCode))
}

// HttpResponseWriter wraps a [http.ResponseWriter] to support
// tracing of the HTTP status code and response size.
type HttpResponseWriter struct {
	span       trace.Span
	wrapped    http.ResponseWriter
	written    int
	statusCode int
}

var httpResponseWriterPool sync.Pool = sync.Pool{
	New: newHttpResponseWriter,
}

// TraceHttpResponse wraps the given [http.ResponseWriter] for tracing purposes
// using the given span. Make sure to always invoke [HttpResponseWriter.Release]
// on the returned [HttpResponseWriter].
func TraceHttpResponse(span trace.Span, w http.ResponseWriter) *HttpResponseWriter {
	return httpResponseWriterPool.Get().(*HttpResponseWriter).init(span, w)
}

func newHttpResponseWriter() any {
	return &HttpResponseWriter{
		statusCode: http.StatusOK,
	}
}

func (w *HttpResponseWriter) init(span trace.Span, wrapped http.ResponseWriter) *HttpResponseWriter {
	w.span = span
	w.wrapped = wrapped
	return w
}

// Release releases this [HttpResponseWriter] and completes tracing
// of the corresponding HTTP request.
func (w *HttpResponseWriter) Release() {
	if w.span.IsRecording() {
		w.span.SetAttributes(semconv.HTTPResponseSize(w.written))
		w.span.SetAttributes(semconv.HTTPResponseStatusCode(w.statusCode))
	}
	w.span = nil
	w.wrapped = nil
	w.written = 0
	w.statusCode = http.StatusOK
	httpResponseWriterPool.Put(w)
}

// Written gets the number written bytes.
func (w *HttpResponseWriter) Written() int {
	return w.written
}

// StatusCode gets the status code.
func (w *HttpResponseWriter) StatusCode() int {
	return w.statusCode
}

// See [http.ResponseWriter.Header]
func (w *HttpResponseWriter) Header() http.Header {
	return w.wrapped.Header()
}

// See [http.ResponseWriter.Write]
func (w *HttpResponseWriter) Write(b []byte) (int, error) {
	written, err := w.wrapped.Write(b)
	if written > 0 {
		w.written += written
	}
	return written, err
}

// See [http.ResponseWriter.WriteHeader]
func (w *HttpResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.wrapped.WriteHeader(statusCode)
}
