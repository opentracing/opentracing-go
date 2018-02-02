/*

Package harness provides a suite of API compatibility checks. They were originally ported from the
OpenTracing Python library's "harness" module.

To run this test suite against your tracer, call harness.RunAPIChecks and provide it a function
that returns a Tracer implementation and a function to call to close it. The function will be
called to create a new tracer before each test in the suite is run, and the returned closer function
will be called after each test is finished.

Several options provide additional checks for your Tracer's behavior: CheckBaggageValues(true)
indicates your tracer supports baggage propagation, CheckExtract(true) tells the suite to test if
the Tracer can extract a trace context from text and binary carriers, and CheckInject(true) tests
if the Tracer can inject the trace context into a carrier.

The UseProbe option provides an APICheckProbe implementation that assists the test suite with
additionally checking if two Spans are part of the same trace, and if a Span and a SpanContext
are part of the same trace. Implementing an APICheckProbe provides additional assertions that
your tracer is working properly.

*/
package harness
