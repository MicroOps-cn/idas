/*
 Copyright © 2023 MicroOps-cn.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package tracing

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/MicroOps-cn/fuck/log"
	uuid "github.com/satori/go.uuid"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type FileTracing struct {
	*stdouttrace.Exporter
	f *os.File
}

func (t *FileTracing) Shutdown(ctx context.Context) error {
	if err := t.f.Sync(); err != nil {
		return err
	}
	return t.f.Close()
}

func NewFileTraceExporter(ctx context.Context, filename string) (sdktrace.SpanExporter, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	exporter, err := stdouttrace.New(
		stdouttrace.WithWriter(f),
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	return &FileTracing{
		Exporter: exporter,
		f:        f,
	}, nil
}

type RetryOptions struct {
	Enabled         bool
	InitialInterval time.Duration
	MaxInterval     time.Duration
	MaxElapsedTime  time.Duration
}

func (c *RetryOptions) UnmarshalJSON(data []byte) (err error) {
	type plain RetryOptions
	*c = RetryOptions{
		Enabled:         true,
		InitialInterval: 5 * time.Second,
		MaxInterval:     30 * time.Second,
		MaxElapsedTime:  time.Minute,
	}
	if string(data) == "false" || string(data) == `"false"` {
		c.Enabled = false
		return
	}
	return json.Unmarshal(data, (*plain)(c))
}

var DefaultRetryConfig = RetryOptions{
	Enabled:         true,
	InitialInterval: 5 * time.Second,
	MaxInterval:     30 * time.Second,
	MaxElapsedTime:  time.Minute,
}

type Compression int

func (c *Compression) UnmarshalJSON(data []byte) (err error) {
	switch string(data) {
	case `"true"`, `true`, `1`, `"1"`, `"gzip"`:
		*c = Compression(otlptracehttp.GzipCompression)
	case `"false"`, `false`, `0`, `"0"`, ``:
		*c = Compression(otlptracehttp.NoCompression)
	}
	return fmt.Errorf("the value can only be one of true, false, or gzip")
}

type TraceOptions struct {
	HTTP        *HTTPClientOptions   `json:"http" yaml:"http" mapstructure:"http"`
	GRPC        *GRPCClientOptions   `json:"grpc" yaml:"grpc" mapstructure:"grpc"`
	Jaeger      *JaegerClientOptions `json:"jaeger" yaml:"jaeger" mapstructure:"jaeger"`
	Zipkin      *ZipkinClientOptions `json:"zipkin" yaml:"zipkin" mapstructure:"zipkin"`
	ServiceName string               `json:"service_name" yaml:"service_name" mapstructure:"service_name"`
}

func (c *TraceOptions) UnmarshalJSON(data []byte) (err error) {
	type plain TraceOptions
	return json.Unmarshal(data, (*plain)(c))
}

type idGenerator struct{}

func (i idGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	tid, err := uuid.FromString(log.GetTraceId(ctx))
	if err != nil {
		tid = uuid.NewV4()
	}

	sid := trace.SpanID{}
	_, _ = rand.Read(sid[:])
	return trace.TraceID(tid), sid
}

func (i idGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	sid := trace.SpanID{}
	_, _ = rand.Read(sid[:])
	return sid
}

var DefaultOptions *TraceOptions

func SetTraceOptions(o *TraceOptions) {
	DefaultOptions = o
}

func NewTraceProvider(ctx context.Context, o *TraceOptions) (p *sdktrace.TracerProvider, err error) {
	var exp sdktrace.SpanExporter
	if o.HTTP != nil {
		exp, err = NewHTTPTraceExporter(ctx, o.HTTP)
	} else if o.GRPC != nil {
		exp, err = NewGRPCTraceExporter(ctx, o.GRPC)
	} else if o.Jaeger != nil {
		exp, err = NewJaegerTraceExporter(ctx, o.Jaeger)
	} else if o.Zipkin != nil {
		exp, err = NewZipkinTraceExporter(ctx, o.Zipkin)
	} else {
		return nil, fmt.Errorf("exporter not specified in tracing configuration")
	}
	if err != nil {
		return nil, err
	}
	if len(o.ServiceName) == 0 {
		o.ServiceName = os.Getenv("APP_NAME")
	}
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(o.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
		sdktrace.WithIDGenerator(&idGenerator{}),
	), nil
}
