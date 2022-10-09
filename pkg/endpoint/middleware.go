/*
 Copyright © 2022 MicroOps-cn.

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

package endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log/level"

	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/logs"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

// InstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
func InstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(method string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			logger := logs.GetContextLogger(ctx)
			defer func(begin time.Time) {
				level.Debug(logger).Log("transport_error", err, "method", method, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func AuthorizationMiddleware(svc service.Service, method string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if needLogin, ok := ctx.Value(global.MetaNeedLogin).(bool); !ok || needLogin {
				if users, ok := ctx.Value(global.MetaUser).([]*models.User); !ok || len(users) == 0 {
					return nil, fmt.Errorf("endpoint authentication failed: system error")
				} else if !svc.Authorization(ctx, users, method) {
					return nil, fmt.Errorf("endpoint authentication failed")
				}
			}
			return next(ctx, request)
		}
	}
}
