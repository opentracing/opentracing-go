# DCP RFC

For background:
- [Main DCP doc](https://paper.dropbox.com/doc/Distributed-Context-Propagation-RGvlvD1NFKYmrJG9vGCES)
- [DCP "Layers" doc](https://paper.dropbox.com/doc/DCP-Layers-and-Semantics-X1Tm1MSiBJECBkWYQKS2I)

## API overview for those adding instrumentation

The higher levels of the opentracing API in golang encourages and takes
advantage of the `context.Context` API. Like it or not, this is what the golang
maintainers want people to do context propagation with inside of a process
(rather than using something like thread-local storage which is unavailable in
golang).

### Initialization

The simplest starting point is `./global.go`. As early as possible, call

    import ".../opentracing"
    
    func main() {
        tracingRecorder := ... // tracing impl specific
        tracingContextIDSource := ... // tracing impl specific
        opentracing.InitGlobalTracer(tracingRecorder, tracingContextIDSource)
        ...
    }

### Creating a Span given an existing Golang `context.Context`

    func xyz(ctx context.Context, ...) {
        ...
        sp, ctx := opentracing.StartSpan("span_name", ctx)
        defer sp.Finish()
		sp.Info("called xyz")
        ...
    }

### Creating a Span given an existing Span

    func xyz(parentSpan opentracing.Span, ...) {
        ...
        sp, ctx := opentracing.StartSpan("span_name", parentSpan)
        defer sp.Finish()
		sp.Info("called xyz")
        ...
    }

### Creating a root Span (i.e., without a known parent)

    func xyz() {
        ...
        sp, ctx := opentracing.StartSpan("span_name")
        defer sp.Finish()
		sp.Info("called xyz")
        ...
    }

### Serializing to the wire

    func makeRequest(req http.Request, ctx context.Context) ... {
        if span := SpanFromGoContext(ctx); span != nil {
            httpClient := &http.Client{}
            httpReq, _ := http.NewRequest("GET", "http://myservice/", nil)
            opentracing.AddContextIDToHttpHeader(span.ContextID(), httpReq.Header)
            resp, err := httpClient.Do(httpReq)
            ...
        }
        ...
    }

### Deerializing from the wire

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		reqCtxID, err := opentracing.GetContextIDFromHttpHeader(
			req.Header, opentracing.GlobalTracer())
		if err != nil {
			panic(err)
		}

		serverSpan, goCtx := opentracing.StartSpan("serverSpan", reqCtxID)
		defer serverSpan.Finish()
        ...
    }

### Note: the singletons are optional

If global singletons make you sad, use `opentracing.NewStandardTracer(...)`
directly and manage ownership of the `opentracing.OpenTracer` explicitly.

## API pointers for those implementing a tracing system

There should be no need for tracing system implementors to worry about the
`opentracing.Span` or `opentracing.OpenTracer` interfaces directly:
`opentracing.NewStandardTracer(...)` should work well enough for most clients.

That said, tracing system authors must provide implementations of:
- `opentracing.ContextID`
- `opentracing.ContextIDSource`
- `opentracing.Recorder`

For a simple working example, see `./dapperish/*.go`.

## TODO items

- An implementation of MultiplexingRecorder (per the comment in `./raw.go`)
  would also make transitions easier for opentracing users.
- There is no sanity-checking or trunctation for large payloads
- There are no safety mechanisms to keep total span memory under a
  user-provided threshold
