// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"net"

	"go.opentelemetry.io/otel/example/grpc/api"
	// "go.opentelemetry.io/otel/example/grpc/config"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/example/grpc/middleware/tracing"
	"go.opentelemetry.io/otel/exporter/trace/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	port = ":7777"
)

// server is used to implement api.HelloServiceServer
type server struct {
	api.UnimplementedHelloServiceServer
}

func initTracer() func() {
	// Create Jaeger Exporter
	exporter, err := jaeger.NewExporter(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "grpc-trace-demo",
			Tags: []core.KeyValue{
				key.String("exporter", "jaeger"),
				key.String("source", "server"),
			},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// For demoing purposes, always sample. In a production application, you should
	// configure this to a trace.ProbabilitySampler set at the desired
	// probability.
	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)
	return func() {
		exporter.Flush()
	}
}

// SayHello implements api.HelloServiceServer
func (s *server) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	log.Printf("Received: %v", in.GetGreeting())
	return &api.HelloResponse{Reply: "Hello " + in.Greeting}, nil
}

func main() {
	// config.Init()
	fn := initTracer()
	defer fn()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(tracing.UnaryServerInterceptor))

	api.RegisterHelloServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
