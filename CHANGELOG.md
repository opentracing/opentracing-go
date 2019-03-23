Changes by Version
==================

1.1.0 (unreleased)
-------------------

- Export StartSpanFromContextWithTracer (#214) <Aaron Delaney>
- Add IsGlobalTracerRegistered() to indicate if a tracer has been registered (#201) <Mike Goldsmith>
- Use Set() instead of Add() in HTTPHeadersCarrier (#191) <jeremyxu2010>
- Apdate license to Apache 2.0 (#181) <Andrea Kao>
- Replace 'golang.org/x/net/context' with 'context' (#176) <Tony Ghita>
- Port of Python opentracing/harness/api_check.py to Go (#146) <chris erway>
- Fix race condition in MockSpan.Context() (#170) <Brad>
- Add PeerHostIPv4.SetString() (#155)  <NeoCN>
- Add a Noop log field type to log to allow for optional fields (#150)  <Matt Ho>


1.0.2 (2017-04-26)
-------------------

- Add more semantic tags (#139) <Rustam Zagirov>


1.0.1 (2017-02-06)
-------------------

- Correct spelling in comments <Ben Sigelman>
- Address race in nextMockID() (#123) <bill fumerola>
- log: avoid panic marshaling nil error (#131) <Anthony Voutas>
- Deprecate InitGlobalTracer in favor of SetGlobalTracer (#128) <Yuri Shkuro>
- Drop Go 1.5 that fails in Travis (#129) <Yuri Shkuro>
- Add convenience methods Key() and Value() to log.Field <Ben Sigelman>
- Add convenience methods to log.Field (2 years, 6 months ago) <Radu Berinde>

1.0.0 (2016-09-26)
-------------------

- This release implements OpenTracing Specification 1.0 (https://opentracing.io/spec)

