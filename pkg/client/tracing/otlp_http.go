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
	"encoding/json"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/MicroOps-cn/idas/pkg/client/internal/tls"
)

func (c *HTTPClientOptions) UnmarshalJSON(data []byte) (err error) {
	type plain HTTPClientOptions
	*c = HTTPClientOptions{
		URLPath:  "/v1/traces",
		Timeout:  time.Second * 10,
		Retry:    DefaultRetryConfig,
		Insecure: true,
	}
	return json.Unmarshal(data, (*plain)(c))
}

type HTTPClientOptions struct {
	Endpoint    string            `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	Timeout     time.Duration     `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Insecure    bool              `json:"insecure" yaml:"insecure" mapstructure:"insecure"`
	Retry       RetryOptions      `json:"retry" yaml:"retry" mapstructure:"retry"`
	Compression Compression       `json:"compression" yaml:"compression" mapstructure:"compression"`
	Header      map[string]string `json:"header" yaml:"header" mapstructure:"header"`
	TLSConfig   *tls.TLSOptions   `json:"tls_config" yaml:"tls_config" mapstructure:"tls_config"`
	URLPath     string            `json:"url_path" yaml:"url_path" mapstructure:"url_path"`
}

func NewHTTPTraceExporter(ctx context.Context, o *HTTPClientOptions) (sdktrace.SpanExporter, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(o.Endpoint),
		otlptracehttp.WithTimeout(o.Timeout),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig(o.Retry)),
		otlptracehttp.WithCompression(otlptracehttp.Compression(o.Compression)),
		otlptracehttp.WithHeaders(o.Header),
		otlptracehttp.WithURLPath(o.URLPath),
	}
	if o.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	} else {
		tlsConfig, err := tls.NewTLSConfig(o.TLSConfig)
		if err != nil {
			return nil, err
		}
		opts = append(opts, otlptracehttp.WithTLSClientConfig(tlsConfig))
	}
	return otlptracehttp.New(ctx, opts...)
}
