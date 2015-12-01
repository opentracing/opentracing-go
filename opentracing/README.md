# DCP RFC

For background:
- [Main DCP doc](https://paper.dropbox.com/doc/Distributed-Context-Propagation-RGvlvD1NFKYmrJG9vGCES)
- [DCP "Layers" doc](https://paper.dropbox.com/doc/DCP-Layers-and-Semantics-X1Tm1MSiBJECBkWYQKS2I)

## API overview for those adding instrumentation

Everyday consumers of this `opentracing` package really only need to worry
about a couple of key abstractions: the `StartSpan` function, the `Span`
interface, and binding a `ProcessRuntime` at `main()`-time. Here are code
snippets demonstrating some important use cases.

#### Singleton initialization

The simplest starting point is `./global.go`. As early as possible, call

    import ".../opentracing"
    
    func main() {
        procRecorder := some_tracing_impl.NewProcessRecorder(...) // tracing impl specific
        traceContextIDSource := some_tracing_impl.NewContextIDSource(...) // tracing impl specific
        opentracing.InitGlobal(procRecorder, traceContextIDSource)
        ...
    }

##### Note: the singletons are optional

If global singletons make you sad, use `opentracing.NewStandardTracer(...)`
directly and manage ownership of the `opentracing.OpenTracer` explicitly.

#### Creating a Span given an existing Span

    func xyz(parentSpan opentracing.Span, ...) {
        ...
        sp := opentracing.JoinTrace("span_name", parentSpan)
        defer sp.Finish()
		sp.Info("called xyz")
        ...
    }

#### Creating a Span given an existing Golang `context.Context`

    func xyz(goCtx context.Context, ...) {
        ...
        sp, goCtx := opentracing.JoinTrace("span_name", goCtx).AddToGoContext(goCtx)
        defer sp.Finish()
		sp.Info("called xyz")
        ...
    }

#### Creating a root Span (i.e., without a known parent)

Just to show that it's not required, we don't call `AddToGoContext` this time.

    func xyz() {
        ...
        sp := opentracing.StartTrace("span_name")
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
            opentracing.AddTraceContextToHttpHeader(span.TraceContext(), httpReq.Header)

            resp, err := httpClient.Do(httpReq)
            ...
        }
        ...
    }

#### Deserializing from the wire

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        // Grab the TraceContext from the HTTP header using the
        // opentracing helper.
        reqTraceCtx, err := opentracing.GetTraceContextFromHttpHeader(
                req.Header, opentracing.GlobalTracer())
        var serverSpan opentracing.Span
        var goCtx context.Context = ...
        if err != nil {
            // Just make a root span.
            serverSpan, goCtx = opentracing.StartSpan("serverSpan").AddToGoContext(goCtx)
        } else {
            // Make a new server-side span that's a child of the span/context sent
            // over the wire.
            serverSpan, goCtx = opentracing.StartSpan("serverSpan", reqTraceCtx).AddToGoContext(goCtx)
        }
        defer serverSpan.Finish()
        ...
    }

## API pointers for those implementing a tracing system

There should be no need for tracing system implementors to worry about the
`opentracing.Span` or `opentracing.OpenTracer` interfaces directly:
`opentracing.NewStandardTracer(...)` should work well enough for most clients.

That said, tracing system authors must provide implementations of:
- `opentracing.TraceContextID`
- `opentracing.TraceContextIDSource`
- `opentracing.ProcessRecorder`

For a small working example, see `./dapperish/*.go`.

## TODO items

- An implementation of MultiplexingRecorder (per the comment in `./raw.go`)
  would also make transitions easier for opentracing users.
- There is no sanity-checking or trunctation for large payloads
- There are no safety mechanisms to keep total span memory under a
  user-provided threshold
