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
	gohttp "net/http"
	"strings"
	"time"

	"github.com/MicroOps-cn/fuck/http"
	"github.com/MicroOps-cn/fuck/log"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log/level"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

// InstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
func InstrumentingMiddleware(duration metrics.Histogram, name string) endpoint.Middleware {
	duration = duration.With("method", name)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			ctx, span := otel.GetTracerProvider().Tracer(config.Get().GetAppName()).Start(ctx, name)
			defer func() {
				if err != nil {
					span.SetStatus(codes.Error, fmt.Sprintf("%+v", err))
				}
				span.End()
			}()
			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(svc service.Service, method string, ps models.Permissions) endpoint.Middleware {
	var postAuditLog func(ctx context.Context, request interface{}, response interface{}, took time.Duration, err error)
	if m := ps.GetMethod(method); m != nil && m.EnableAudit {
		postAuditLog = func(ctx context.Context, request interface{}, response interface{}, took time.Duration, err error) {
			logger := log.GetContextLogger(ctx)
			defer func() {
				if r := recover(); r != nil {
					level.Error(logger).Log("msg", "failed to post event log", "err", err)
				}
			}()
			var stdReq *gohttp.Request
			if requester, ok := request.(RestfulRequester); ok {
				stdReq = requester.GetRestfulRequest().Request
			} else if restReq, ok := request.(restful.Request); ok {
				stdReq = restReq.Request
			} else {
				stdReq = &gohttp.Request{}
			}
			remoteAddr := http.GetRemoteAddr(stdReq, config.Get().Security.TrustIp)
			var username, userId string
			if user, ok := ctx.Value(global.MetaUser).(*models.User); ok || user != nil {
				username, userId = user.Username, user.Id
			}
			call := ps.Get(method)
			var msg string
			var status bool
			if err != nil {
				msg = fmt.Sprintf("Failed to call %s (%s): %s", call[0].Name, call[0].Description, err)
			} else if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
				msg = fmt.Sprintf("Failed to call %s (%s): %s", call[0].Name, call[0].Description, f.Failed())
			} else {
				msg = fmt.Sprintf("Success to call %s (%s)", call[0].Name, call[0].Description)
				status = true
			}
			l := map[string]string{
				"Title":           msg,
				"RequestLine":     fmt.Sprintf("%s %s %s", stdReq.Method, stdReq.URL.String(), stdReq.Proto),
				"X-Forwarded-For": strings.Join(stdReq.Header.Values("X-Forwarded-For"), ","),
				"RemoteAddr":      stdReq.RemoteAddr,
				"User-Agent":      strings.Join(stdReq.Header.Values("User-Agent"), ","),
				//"Request":         w.JSONStringer(reqData).String(),
				"Type": "call_endpoint",
			}
			if e := svc.PostEventLog(ctx, log.GetTraceId(ctx), userId, username, remoteAddr, method, msg, status, took, l); e != nil {
				level.Error(logger).Log("failed to post event log", "err", e)
			}
		}
	}
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			logger := log.GetContextLogger(ctx)
			level.Debug(logger).Log("msg", "call method", "method", method)
			defer func(begin time.Time) {
				took := time.Since(begin)
				if postAuditLog != nil {
					postAuditLog(ctx, request, response, took, err)
				}
				if err != nil {
					level.Debug(logger).Log("msg", "method call finished", "transport_error", err, "method", method, "took", took)
				} else {
					level.Debug(logger).Log("msg", "method call finished", "method", method, "took", took)
				}
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func AuthorizationMiddleware(svc service.Service, method string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if needLogin, ok := ctx.Value(global.MetaNeedLogin).(bool); !ok || needLogin {
				if user, ok := ctx.Value(global.MetaUser).(*models.User); !ok || user == nil {
					return nil, errors.NewServerError(401, "need login")
				} else if !svc.Authorization(ctx, user, method) {
					if forceOk, ok := ctx.Value(global.MetaForceOk).(bool); ok && forceOk {
						return nil, nil
					}
					return nil, errors.NewServerError(403, "forbidden")
				}
			}
			return next(ctx, request)
		}
	}
}
