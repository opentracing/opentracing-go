# OpenTracing API for Go

This package is a Go platform API for OpenTracing.

## Required Reading

In order to understand the Go platform API, one must first be familiar with the
[OpenTracing project](http://opentracing.io) and
[terminology](http://opentracing.io/spec/) more generally.

## API overview for those adding instrumentation

Everyday consumers of this `opentracing` package really only need to worry
about a couple of key abstractions: the `StartTrace` function, the `Span`
interface, and binding a `Tracer` at `main()`-time. Here are code snippets
demonstrating some important use cases.

#### Singleton initialization

The simplest starting point is `./default_tracer.go`. As early as possible, call

```
    import ".../opentracing-go"
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

#### Starting an empty trace by creating a "root span"

```
    func xyz() {
        ...
        sp := opentracing.StartTrace("span_name")
        defer sp.Finish()
        sp.Info("called xyz")
        ...
    }
```

#### Creating a Span given an existing Span

```
    func xyz(parentSpan opentracing.Span, ...) {
        ...
        sp := opentracing.JoinTrace("span_name", parentSpan)
        defer sp.Finish()
        sp.Info("called xyz")
        ...
    }
```

#### Creating a Span given an existing Golang `context.Context`

Additionally, this example demonstrates how to get a `context.Context`
associated with any `opentracing.Span` instance.

```
    func xyz(goCtx context.Context, ...) {
        ...
        sp, goCtx := opentracing.JoinTrace("span_name", goCtx).AddToGoContext(goCtx)
        defer sp.Finish()
        sp.Info("called xyz")
        ...
    }
```

#### Serializing to the wire

```
    func makeSomeRequest(ctx context.Context) ... {
        if span := opentracing.SpanFromGoContext(ctx); span != nil {
            httpClient := &http.Client{}
            httpReq, _ := http.NewRequest("GET", "http://myservice/", nil)

            // Transmit the span's TraceContext as an HTTP header on our
            // outbound request.
            opentracing.AddTraceContextToHeader(
                span.TraceContext(),
                httpReq.Header,
                opentracing.DefaultTracer())

            resp, err := httpClient.Do(httpReq)
            ...
        }
        ...
    }
```

#### Deserializing from the wire

```
    http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        // Grab the TraceContext from the HTTP header using the
        // opentracing helper.
        reqTraceCtx, err := opentracing.TraceContextFromHeader(
                req.Header, opentracing.GlobalTracer())
        var serverSpan opentracing.Span
        var goCtx context.Context = ...
        if err != nil {
            // Just make a root span.
            serverSpan, goCtx = opentracing.StartTrace("serverSpan").AddToGoContext(goCtx)
        } else {
            // Make a new server-side span that's a child of the span/context sent
            // over the wire.
            serverSpan, goCtx = opentracing.JoinTrace(
                "serverSpan", reqTraceCtx).AddToGoContext(goCtx)
        }
        defer serverSpan.Finish()
        ...
    }
```

#### Goroutine-safety

The entire public API is goroutine-safe and does not require external
synchronization.

## API pointers for those implementing a tracing system

There should be no need for most tracing system implementors to worry about the
`opentracing.Span` or `opentracing.Tracer` interfaces directly:
`standardtracer.New(...)` should work well enough in most circumstances.

In order to integrate with `standardtracer`, tracing system authors are
expected to provide implementations of:
- `opentracing.TraceContext`
- `opentracing.TraceContextSource`
- `standardtracer.Recorder`

For a small working example, see `../examples/dapperish/*.go`.
