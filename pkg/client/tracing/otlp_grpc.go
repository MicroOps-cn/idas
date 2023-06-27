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
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"

	"github.com/MicroOps-cn/idas/pkg/client/internal/tls"
)

type GRPCClientOptions struct {
	Endpoint           string            `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	Timeout            time.Duration     `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Insecure           bool              `json:"insecure" yaml:"insecure" mapstructure:"insecure"`
	Retry              RetryOptions      `json:"retry" yaml:"retry" mapstructure:"retry"`
	Compression        Compression       `json:"compression" yaml:"compression" mapstructure:"compression"`
	Header             map[string]string `json:"header" yaml:"header" mapstructure:"header"`
	TLSConfig          tls.TLSOptions    `json:"tls_config" yaml:"tls_config" mapstructure:"tls_config"`
	URLPath            string            `json:"url_path" yaml:"url_path" mapstructure:"url_path"`
	ReconnectionPeriod time.Duration     `json:"reconnection_period" yaml:"reconnection_period" mapstructure:"reconnection_period"`
	ServiceConfig      string            `json:"service_config" yaml:"service_config" mapstructure:"service_config"`
}

func NewGRPCTraceExporter(ctx context.Context, o *GRPCClientOptions) (sdktrace.SpanExporter, error) {
	tlsConfig, err := tls.NewTLSConfig(&o.TLSConfig)
	if err != nil {
		return nil, err
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(o.Endpoint),
		otlptracegrpc.WithTimeout(o.Timeout),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig(o.Retry)),
		otlptracegrpc.WithHeaders(o.Header),
		otlptracegrpc.WithServiceConfig(o.ServiceConfig),
		otlptracegrpc.WithReconnectionPeriod(o.ReconnectionPeriod),
		otlptracegrpc.WithTLSCredentials(credentials.NewTLS(tlsConfig)),
	}
	if o.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	if o.Compression == 1 {
		opts = append(opts, otlptracegrpc.WithCompressor(gzip.Name))
	}
	return otlptracegrpc.New(ctx, opts...)
}
