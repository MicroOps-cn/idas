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
	"strconv"
	"time"

	"github.com/MicroOps-cn/fuck/log"
	w "github.com/MicroOps-cn/fuck/wrapper"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2868"
	"layeh.com/radius/rfc2869"

	"github.com/MicroOps-cn/idas/pkg/endpoint"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/logs"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

type SecretSourceFunc func(r *radius.Request) ([]byte, error)

func (f SecretSourceFunc) RADIUSSecret(r *radius.Request) ([]byte, error) {
	return f(r)
}

type ResponseWriter struct {
	radius.ResponseWriter
	code radius.Code
}

func (r *ResponseWriter) Write(packet *radius.Packet) error {
	r.code = packet.Code
	return r.ResponseWriter.Write(packet)
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
		wr := &ResponseWriter{ResponseWriter: writer}
		defer func() {
			if r := recover(); r != nil {
				if req.Packet == nil {
					_ = wr.Write(req.Response(radius.CodeAccessReject))
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
				logs.WrapKeyName("code"), wr.code,
			)
			level.Info(logger).Log(logs.WrapKeyName("totalTime"), fmt.Sprintf("%dms", time.Since(start).Milliseconds()))
		}()
		var err error
		raw := req.GetRaw()
		req.Packet = &radius.Packet{
			Code:       radius.Code(raw[0]),
			Identifier: raw[1],
		}
		req.Attributes, err = radius.ParseAttributes(raw[20:])
		if err != nil {
			level.Error(logger).Log("msg", "failed to parse request body", "err", err)
			_ = wr.Write(req.Response(radius.CodeAccessReject))
			return
		}
		headers := map[string]string{}
		var kvs []interface{}
		for _, attribute := range req.Attributes {
			name, ok := RadiusTypeName[attribute.Type]
			if !ok {
				name = strconv.Itoa(int(attribute.Type))
			}
			if attribute.Type == rfc2865.UserName_Type {
				kvs = append(kvs, log.WrapKeyName(name), radius.String(attribute.Attribute))
			} else if attribute.Type == rfc2865.NASIdentifier_Type {
				kvs = append(kvs, log.WrapKeyName(name), radius.String(attribute.Attribute))
			} else if attribute.Type == rfc2865.NASIPAddress_Type {
				addr, err := radius.IPAddr(attribute.Attribute)
				if err != nil {
					kvs = append(kvs, log.WrapKeyName(name), fmt.Sprintf("[base64]%s", base64.StdEncoding.EncodeToString(attribute.Attribute)))
				} else {
					kvs = append(kvs, log.WrapKeyName(name), addr.String())
				}
			} else if attribute.Type == rfc2865.UserPassword_Type ||
				attribute.Type == rfc2865.CHAPPassword_Type ||
				attribute.Type == rfc2865.ReplyMessage_Type ||
				attribute.Type == rfc2868.TunnelPassword_Type ||
				attribute.Type == rfc2869.ARAPPassword_Type {
				kvs = append(kvs, log.WrapKeyName(name), "*****************")
			} else {
				headers[name] = fmt.Sprintf("[base64]%s", base64.StdEncoding.EncodeToString(attribute.Attribute))
			}
		}
		level.Info(kitlog.With(logger, kvs...)).Log(
			"msg", "HTTP request received.",
			logs.TitleKey, "request",
			logs.WrapKeyName("remoteAddr"), req.RemoteAddr,
			logs.WrapKeyName("code"), req.Code,
			log.WrapKeyName("header"), w.JSONStringer(headers),
			//logs.WrapKeyName("body"), w.NewStringer(func() string {
			//	return base64.StdEncoding.EncodeToString(req.GetRaw())
			//}),
		)
		fc.ProcessFilter(wr, req)
	}
}

func RadiusAppFilter(endpoints endpoint.Set) radius.FilterFunction {
	return func(w radius.ResponseWriter, r *radius.Request, fc *radius.FilterChain) {
		ctx := r.Context()
		logger := log.GetContextLogger(ctx)
		username, _ := r.Attributes.Lookup(rfc2865.UserName_Type)
		id, _ := r.Attributes.Lookup(rfc2865.NASIdentifier_Type)
		if s, err := endpoints.GetAppAndKeyFromKeyId(r.Context(), &HTTPRequest[endpoint.GetAppKeyRequestData]{Data: endpoint.GetAppKeyRequestData{
			Username: string(username),
			Key:      string(id),
		}}); err != nil {
			level.Error(logger).Log("msg", "failed to get app and key by key id", "err", err)
			w.Write(r.Response(radius.CodeAccessReject))
		} else if appKey, ok := s.(*endpoint.GetAppKeyResponseData); ok && appKey != nil {
			if appKey.App.GrantType&models.AppMeta_radius > 0 && (len(appKey.App.Users) > 0) {
				pkt, err := radius.Parse(r.GetRaw(), []byte(appKey.Key.Secret))
				if err != nil {
					level.Error(logger).Log("msg", "failed to parse request body", "err", err)
					w.Write(r.Response(radius.CodeAccessReject))
					return
				}
				r.Packet = pkt
				ctx = context.WithValue(ctx, global.MetaAppSecretHash, appKey.Key.Secret)
				r.Secret = []byte(appKey.Key.Secret)
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

func RadiusValidateAuthRequestFilter(w radius.ResponseWriter, r *radius.Request, fc *radius.FilterChain) {
	logger := log.GetContextLogger(r.Context())
	if _, ok := r.Lookup(rfc2865.UserName_Type); !ok {
		level.Error(logger).Log("msg", "invalid request", "err", "UserName attribute missing")
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}
	if _, ok := r.Lookup(rfc2865.NASIdentifier_Type); !ok {
		level.Error(logger).Log("msg", "invalid request", "err", "NASIdentifier attribute missing")
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}
	fc.ProcessFilter(w, r)
}

type HandlerChain []func(w radius.ResponseWriter, r *radius.Request, ch HandlerChain)

func (h HandlerChain) ServeRADIUS(w radius.ResponseWriter, r *radius.Request) {
	if len(h) == 0 {
		return
	}
	h[0](w, r, h[1:])
}

func NewRadiusService(ctx context.Context, endpoints endpoint.Set) *radius.PacketServer {
	server := radius.PacketServer{
		Filters: []radius.FilterFunction{
			RadiusLoggingFilter(ctx),
			RadiusValidateAuthRequestFilter,
			RadiusAppFilter(endpoints),
			RadiusEAPAuthFilter,
		},
		Handler: HandlerChain{
			RadiusEAPHandler,
			func(w radius.ResponseWriter, r *radius.Request, ch HandlerChain) {
				username := rfc2865.UserName_GetString(r.Packet)
				password, err := rfc2865.UserPassword_LookupString(r.Packet)
				if err == radius.ErrNoAttribute {
					level.Debug(log.GetContextLogger(r.Context())).Log("msg", "invalid request", "err", "UserPassword attribute missing")
					w.Write(r.Response(radius.CodeAccessReject))
					return
				}
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
			},
		},
		SecretSource: SecretSourceFunc(func(r *radius.Request) (secret []byte, err error) {
			return []byte(r.Context().Value(global.MetaAppSecretHash).(string)), nil
		}),
	}
	return &server
}
