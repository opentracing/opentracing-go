[![Gitter chat](http://img.shields.io/badge/gitter-join%20chat%20%E2%86%92-brightgreen.svg)](https://gitter.im/opentracing/public) [![Build Status](https://travis-ci.org/opentracing/opentracing-go.svg?branch=master)](https://travis-ci.org/opentracing/opentracing-go) [![GoDoc](https://godoc.org/github.com/opentracing/opentracing-go?status.svg)](http://godoc.org/github.com/opentracing/opentracing-go)

# OpenTracing API for Go

This package is a Go platform API for OpenTracing.

## Required Reading

In order to understand the Go platform API, one must first be familiar with the
[OpenTracing project](http://opentracing.io) and
[terminology](http://opentracing.io/spec/) more generally.

## API overview for those adding instrumentation

Everyday consumers of this `opentracing` package really only need to worry
about a couple of key abstractions: the `StartSpan` function, the `Span`
interface, and binding a `Tracer` at `main()`-time. Here are code snippets
demonstrating some important use cases.

#### Singleton initialization

The simplest starting point is `./default_tracer.go`. As early as possible, call

```go
    import "github.com/opentracing/opentracing-go"
    import ".../some_tracing_impl"

    func main() {
        tracerImpl := some_tracing_impl.New(...) // tracing impl specific
        opentracing.InitGlobalTracer(tracerImpl)
        ...
    }
```

##### Non-Singleton initialization

If you prefer direct control to singletons, manage ownership of the
`opentracing.Tracer` implementation explicitly.

#### Creating a Span given an existing Golang `context.Context`

If you use `context.Context` in your application, OpenTracing's Go library will
happily use it for `Span` propagation. To start a new (child) `Span`, you can use
`StartSpanFromContext`.

```go
    func xyz(ctx context.Context, ...) {
        ...
        span, ctx := opentracing.StartSpanFromContext(ctx, "operation_name")
        defer span.Finish()
        span.LogEvent("xyz_called")
        ...
    }
```

#### Starting an empty trace by creating a "root span"

It's always possible to create a "root" (parentless) `Span`.

```go
    func xyz() {
        ...
        sp := opentracing.StartSpan("operation_name")
        defer sp.Finish()
        sp.LogEvent("xyz_called")
        ...
    }
```

#### Creating a (child) Span given an existing (parent) Span

```go
    func xyz(parentSpan opentracing.Span, ...) {
        ...
        sp := opentracing.StartChildSpan(parentSpan, "operation_name")
        defer sp.Finish()
        sp.LogEvent("xyz_called")
        ...
    }
```

#### Serializing to the wire

```go
    func makeSomeRequest(ctx context.Context) ... {
        if span := opentracing.SpanFromContext(ctx); span != nil {
            httpClient := &http.Client{}
            httpReq, _ := http.NewRequest("GET", "http://myservice/", nil)

            // Transmit the span's TraceContext as HTTP headers on our
            // outbound request.
            tracer.Inject(
                span,
                opentracing.TextMap,
                opentracing.HTTPHeaderTextMapCarrier(httpReq.Header))

            resp, err := httpClient.Do(httpReq)
            ...
        }
        ...
    }
```

#### Deserializing from the wire

```go
    http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        serverSpan, err := opentracing.GlobalTracer().Join(
            "serverSpan",
            opentracing.TextMap,
            opentracing.HTTPHeaderTextMapCarrier(req.Header))

        if err != nil {
            // Create a root span if necessary
            serverSpan = opentracing.StartTrace("serverSpan")
        }
        var goCtx context.Context = ...
        goCtx, _ = opentracing.ContextWithSpan(goCtx, serverSpan)
        defer serverSpan.Finish()
        ...
    }
```

#### Goroutine-safety

The entire public API is goroutine-safe and does not require external
synchronization.

## API pointers for those implementing a tracing system

Tracing system implementors may be able to reuse or copy-paste-modify the `basictracer` package, found [here](https://github.com/opentracing/basictracer-go). In particular, see `basictracer.New(...)`.

## API compatibility

For the time being, "mild" backwards-incompatible changes may be made without changing the major version number. As OpenTracing and `opentracing-go` mature, backwards compatibility will become more of a priority.
