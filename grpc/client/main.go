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
	"time"

	"go.opentelemetry.io/otel/example/grpc/api"
	// "go.opentelemetry.io/otel/example/grpc/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"

	"go.opentelemetry.io/otel/example/grpc/middleware/tracing"
	"go.opentelemetry.io/otel/exporter/trace/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() func() {
	// Create Jaeger Exporter
	exporter, err := jaeger.NewExporter(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "grpc-trace-demo",
			Tags: []core.KeyValue{
				key.String("exporter", "jaeger"),
				key.String("source", "client"),
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

func main() {
	// config.Init()
	fn := initTracer()
	defer fn()

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":7777", grpc.WithInsecure(), grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer func() { _ = conn.Close() }()

	c := api.NewHelloServiceClient(conn)

	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-api-client-us-east-1",
		"user-id", "some-test-user-id",
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	response, err := c.SayHello(ctx, &api.HelloRequest{Greeting: "World"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Reply)
}
