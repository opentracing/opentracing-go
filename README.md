# Golang OpenTracing API

For background:
- [Main DCP doc](https://paper.dropbox.com/doc/Distributed-Context-Propagation-RGvlvD1NFKYmrJG9vGCES)
- [DCP "Layers" doc](https://paper.dropbox.com/doc/DCP-Layers-and-Semantics-X1Tm1MSiBJECBkWYQKS2I)

## API overview for those adding instrumentation

Everyday consumers of this `opentracing` package really only need to worry
about a couple of key abstractions: the `StartSpan` function, the `Span`
interface, and binding a `Recorder` at `main()`-time. Here are code snippets
demonstrating some important use cases.

#### Singleton initialization

The simplest starting point is `./default_tracer.go`. As early as possible, call

    import ".../opentracing-go"
    
    func main() {
        procRecorder := some_tracing_impl.NewRecorder(...) // tracing impl specific
        traceContextSource := some_tracing_impl.NewContextSource(...) // tracing impl specific
        opentracing.InitDefaultTracer(procRecorder, traceContextSource)
        ...
    }

##### Non-Singleton initialization

If you prefer direct control to singletons, use `standardtracer.New(...)`
directly and manage ownership of the `opentracing.Tracer` implementation
explicitly.

#### Starting an empty trace by creating a "root span"

    func xyz() {
        ...
        sp := opentracing.StartTrace("span_name")
        defer sp.Finish()
		sp.Info("called xyz")
        ...
    }

#### Creating a Span given an existing Span

    func xyz(parentSpan opentracing.Span, ...) {
        ...
        sp := opentracing.JoinTrace("span_name", parentSpan)
        defer sp.Finish()
		sp.Info("called xyz")
        ...
    }

#### Creating a Span given an existing Golang `context.Context`

Additionally, this example demonstrates how to get a `context.Context`
associated with any `opentracing.Span` instance.

    func xyz(goCtx context.Context, ...) {
        ...
        sp, goCtx := opentracing.JoinTrace("span_name", goCtx).AddToGoContext(goCtx)
        defer sp.Finish()
		sp.Info("called xyz")
        ...
    }

#### Serializing to the wire

    func makeSomeRequest(ctx context.Context) ... {
        if span := SpanFromGoContext(ctx); span != nil {
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

#### Deserializing from the wire

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
            serverSpan, goCtx = opentracing.JoinTrace("serverSpan", reqTraceCtx).AddToGoContext(goCtx)
        }
        defer serverSpan.Finish()
        ...
    }

#### Goroutine-safety

The entire public API is goroutine-safe and does not require external
synchronization.

## API pointers for those implementing a tracing system

There should be no need for most tracing system implementors to worry about the
`opentracing.Span` or `opentracing.Tracer` interfaces directly:
`standardtracer.New(...)` should work well enough in most circumstances.

That said, tracing system authors are expected to provide implementations of:
- `opentracing.TraceContext`
- `opentracing.TraceContextSource`
- `standardtracer.Recorder`

For a small working example, see `../examples/dapperish/*.go`.

## TODO items

- An implementation of MultiplexingRecorder (per the comment in `./raw.go`)
  would also make transitions easier for opentracing users.
- There is no sanity-checking or trunctation for large payloads
- There are no safety mechanisms to keep total span memory under a
  user-provided threshold
