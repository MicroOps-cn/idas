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

	"go.opentelemetry.io/otel/exporters/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type ZipkinClientOptions struct {
	Endpoint string `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
}

func NewZipkinTraceExporter(_ context.Context, o *ZipkinClientOptions) (sdktrace.SpanExporter, error) {
	return zipkin.New(o.Endpoint)
}
