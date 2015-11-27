# DCP RFC

For background:
- [Main DCP doc](https://paper.dropbox.com/doc/Distributed-Context-Propagation-RGvlvD1NFKYmrJG9vGCES)
- [DCP "Layers" doc](https://paper.dropbox.com/doc/DCP-Layers-and-Semantics-X1Tm1MSiBJECBkWYQKS2I)

## API overview for those adding instrumentation

Everyday consumers of this `opentracing` package really only need to worry
about a couple of key abstractions: the `StartSpan` function, the `Span`
interface, and golang's `context.Context` idiom. Here are code snippets
demonstrating some important use cases.

A note about `context.Context`: the opentracing API here encourages (but does
not require) the use of golang-team's `context.Context` scheme. Whether we like
it or not (this author is on the fence), `context.Context` is here to stay,
goroutine-local-storage bedamned. When in Rome, do as the Romans do (etc).

#### Singleton initialization

The simplest starting point is `./global.go`. As early as possible, call

    import ".../opentracing"
    
    func main() {
        tracingRecorder := ... // tracing impl specific
        tracingContextIDSource := ... // tracing impl specific
        opentracing.InitGlobalTracer(tracingRecorder, tracingContextIDSource)
        ...
    }

##### Note: the singletons are optional

If global singletons make you sad, use `opentracing.NewStandardTracer(...)`
directly and manage ownership of the `opentracing.OpenTracer` explicitly.

#### Creating a Span given an existing Golang `context.Context`

    func xyz(ctx context.Context, ...) {
        ...
        sp, ctx := opentracing.StartSpan("span_name", ctx)
        defer sp.Finish()
		sp.Info("called xyz")
        ...
    }

#### Creating a Span given an existing Span

    func xyz(parentSpan opentracing.Span, ...) {
        ...
        sp, ctx := opentracing.StartSpan("span_name", parentSpan)
        defer sp.Finish()
		sp.Info("called xyz")
        ...
    }

#### Creating a root Span (i.e., without a known parent)

    func xyz() {
        ...
        sp, ctx := opentracing.StartSpan("span_name")
        defer sp.Finish()
		sp.Info("called xyz")
        ...
    }

#### Serializing to the wire

    func makeSomeRequest(ctx context.Context) ... {
        if span := SpanFromGoContext(ctx); span != nil {
            httpClient := &http.Client{}
            httpReq, _ := http.NewRequest("GET", "http://myservice/", nil)

			// Transmit the span's ContextID as an HTTP header on our outbound
            // request.
            opentracing.AddContextIDToHttpHeader(span.ContextID(), httpReq.Header)

            resp, err := httpClient.Do(httpReq)
            ...
        }
        ...
    }

#### Deserializing from the wire

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// Grab the ContextID from the HTTP header using the
                // opentracing helper.
                reqCtxID, err := opentracing.GetContextIDFromHttpHeader(
                        req.Header, opentracing.GlobalTracer())
                var serverSpan opentracing.Span
                var goCtx context.Context
                if err != nil {
                    // Just make a root span.
                    serverSpan, goCtx = opentracing.StartSpan("serverSpan")
                } else {
                    // Make a new server-side span that's a child of the span/context sent
                    // over the wire.
                    serverSpan, goCtx = opentracing.StartSpan("serverSpan", reqCtxID)
                }
		defer serverSpan.Finish()
        ...
    }

## API pointers for those implementing a tracing system

There should be no need for tracing system implementors to worry about the
`opentracing.Span` or `opentracing.OpenTracer` interfaces directly:
`opentracing.NewStandardTracer(...)` should work well enough for most clients.

That said, tracing system authors must provide implementations of:
- `opentracing.ContextID`
- `opentracing.ContextIDSource`
- `opentracing.ComponentRecorder`

For a simple working example, see `./dapperish/*.go`.

## TODO items

- An implementation of MultiplexingRecorder (per the comment in `./raw.go`)
  would also make transitions easier for opentracing users.
- There is no sanity-checking or trunctation for large payloads
- There are no safety mechanisms to keep total span memory under a
  user-provided threshold
