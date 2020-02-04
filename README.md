# Opentelemetry gRPC example with jaeger implementation
The [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go) repo has separate [jaeger](https://github.com/open-telemetry/opentelemetry-go/tree/master/example/jaeger) and [gRPC](https://github.com/open-telemetry/opentelemetry-go/tree/master/example/grpc) examples. This is an experiment in combining the two, since there was no obvious jaeger Inject Extract example.

The project uses most of the gRPC example however the config folder was omitted due to a the jaeger tracer implementation.

This was done as an experiment for possible implementation, the opentelemetry is [still being developed](https://opentelemetry.io/project-status/) at the time of pushing this, so it shouldn't be used in a production tracing.

Keep checking [opentelemetry](https://opentelemetry.io/) for any new updates :)
