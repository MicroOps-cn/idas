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

package transport

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"runtime"
	"time"

	"github.com/MicroOps-cn/fuck/log"
	w "github.com/MicroOps-cn/fuck/wrapper"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"

	"github.com/MicroOps-cn/idas/pkg/endpoint"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/logs"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

type SecretSourceFunc func(r *radius.Request) ([]byte, error)

func (f SecretSourceFunc) RADIUSSecret(r *radius.Request) ([]byte, error) {
	return f(r)
}

func RadiusLoggingFilter(pctx context.Context) radius.FilterFunction {
	return func(writer radius.ResponseWriter, req *radius.Request, fc *radius.FilterChain) {
		ctx := req.Context()
		if ctx == nil {
			ctx = pctx
		}
		traceId := log.NewTraceId()
		var logger kitlog.Logger
		ctx, logger = log.NewContextLogger(ctx, log.WithTraceId(traceId))
		req = req.WithContext(ctx)
		start := time.Now()
		defer func() {
			if r := recover(); r != nil {
				if req.Packet == nil {
					writer.Write(req.Response(radius.CodeAccessReject))
					return
				}
				buf := bytes.NewBufferString(fmt.Sprintf("recover from panic situation: - %v\n", r))
				for i := 2; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					buf.WriteString(fmt.Sprintf("    %s:%d\n", file, line))
				}
				level.Error(logger).Log("msg", buf.String())
			}
			logger = kitlog.With(logger,
				"msg", "HTTP response send.",
				logs.TitleKey, "response",
			)
			level.Info(logger).Log(logs.WrapKeyName("totalTime"), fmt.Sprintf("%dms", time.Since(start).Milliseconds()))
		}()

		level.Info(logger).Log(
			"msg", "HTTP request received.",
			logs.TitleKey, "request",
			logs.WrapKeyName("remoteAddr"), req.RemoteAddr,
			logs.WrapKeyName("code"), req.Code,
			logs.WrapKeyName("body"), w.NewStringer(func() string {
				return base64.StdEncoding.EncodeToString(req.GetRaw())
			}),
		)
		fc.ProcessFilter(writer, req)
	}
}

func RadiusAppFilter(endpoints endpoint.Set) radius.FilterFunction {
	return func(w radius.ResponseWriter, r *radius.Request, fc *radius.FilterChain) {
		ctx := r.Context()
		logger := log.GetContextLogger(ctx)
		if attr, err := radius.ParseAttributes(r.GetRaw()[20:]); err != nil {
			level.Error(logger).Log("msg", "failed to parse request body", "err", err)
			w.Write(r.Response(radius.CodeAccessReject))
			return
		} else if id, ok := attr.Lookup(rfc2865.NASIdentifier_Type); !ok {
			level.Error(logger).Log("msg", "failed to get app id from client request", "err", err)
			w.Write(r.Response(radius.CodeAccessReject))
			return
		} else if username, ok := attr.Lookup(rfc2865.UserName_Type); !ok {
			level.Error(logger).Log("msg", "failed to get username from client request", "err", err)
			w.Write(r.Response(radius.CodeAccessReject))
			return
		} else {
			if s, err := endpoints.GetAppAndKeyFromKeyId(r.Context(), &HTTPRequest[endpoint.GetAppKeyRequestData]{Data: endpoint.GetAppKeyRequestData{
				Username: string(username),
				Key:      string(id),
			}}); err != nil {
				level.Error(logger).Log("msg", "failed to get app and key by key id", "err", err)
				w.Write(r.Response(radius.CodeAccessReject))
			} else if appKey, ok := s.(*endpoint.GetAppKeyResponseData); ok && appKey != nil {
				ctx = context.WithValue(ctx, global.MetaAppSecretHash, appKey.Key.Secret)
				r.Secret = []byte(appKey.Key.Secret)
				if appKey.App.GrantType&models.AppMeta_radius > 0 && (len(appKey.App.Users) > 0) {
					r = r.WithContext(ctx)
					fc.ProcessFilter(w, r)
				} else {
					w.Write(r.Response(radius.CodeAccessReject))
				}
			} else {
				w.Write(r.Response(radius.CodeAccessReject))
			}
		}
	}
}

func NewRadiusService(ctx context.Context, endpoints endpoint.Set) *radius.PacketServer {
	server := radius.PacketServer{
		Filters: []radius.FilterFunction{RadiusLoggingFilter(ctx), RadiusAppFilter(endpoints)},
		Handler: radius.HandlerFunc(func(w radius.ResponseWriter, r *radius.Request) {
			username := rfc2865.UserName_GetString(r.Packet)
			password := rfc2865.UserPassword_GetString(r.Packet)
			var code radius.Code
			authReq := &HTTPRequest[endpoint.AuthenticationRequest]{
				Data: endpoint.AuthenticationRequest{
					AuthKey:    username,
					AuthSecret: password,
				},
			}
			if u, err := endpoints.Authentication(r.Context(), authReq); err != nil {
				code = radius.CodeAccessReject
			} else if user, ok := u.(*models.User); ok && len(user.Id) > 0 {
				code = radius.CodeAccessAccept
			} else {
				code = radius.CodeAccessReject
			}
			w.Write(r.Response(code))
		}),
		SecretSource: SecretSourceFunc(func(r *radius.Request) (secret []byte, err error) {
			return []byte(r.Context().Value(global.MetaAppSecretHash).(string)), nil
		}),
	}
	return &server
}
