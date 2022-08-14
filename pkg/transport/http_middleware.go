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

package transport

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/MicroOps-cn/idas/pkg/endpoint"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/logs"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/utils/buffer"
	"github.com/MicroOps-cn/idas/pkg/utils/httputil"
	"github.com/MicroOps-cn/idas/pkg/utils/sign"
	w "github.com/MicroOps-cn/idas/pkg/utils/wrapper"
)

func HTTPAuthenticationFilter(endpoints endpoint.Set) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, filterChan *restful.FilterChain) {
		if req.SelectedRoute() == nil {
			errorEncoder(req.Request.Context(), errors.NewServerError(http.StatusNotFound, "Not Found: "+req.Request.RequestURI), resp)
			return
		}
		if needLogin, ok := req.SelectedRoute().Metadata()[global.MetaNeedLogin].(bool); ok {
			//
			req.Request = req.Request.WithContext(context.WithValue(req.Request.Context(), global.MetaNeedLogin, needLogin))
			if !needLogin {
				filterChan.ProcessFilter(req, resp)
				return
			}
		}

		errorHandler := errorEncoder

		loginSessionID, err := req.Request.Cookie(global.LoginSession)
		if err == nil {
			req.Request = req.Request.WithContext(context.WithValue(req.Request.Context(), global.LoginSession, loginSessionID.Value))
			if user, err := endpoints.GetSessionByToken(req.Request.Context(), endpoint.GetSessionParams{
				Token:     loginSessionID.Value,
				TokenType: models.TokenTypeLoginSession,
			}); err == nil {
				if len(user.([]*models.User)) >= 0 {
					req.Request = req.Request.WithContext(context.WithValue(req.Request.Context(), global.MetaUser, user))
					filterChan.ProcessFilter(req, resp)
					return
				}
			}
		}
		authReq := HTTPRequest[endpoint.AuthenticationRequest]{}
		if username, password, ok := req.Request.BasicAuth(); ok {
			authReq.Data.AuthKey = username
			authReq.Data.AuthSecret = password
		} else if auth := req.Request.Header.Get("Authorization"); len(auth) != 0 {
			if strings.HasPrefix(auth, "Bearer ") {
				if user, err := endpoints.GetSessionByToken(req.Request.Context(), endpoint.GetSessionParams{
					Token:     strings.TrimPrefix(auth, "Bearer "),
					TokenType: models.TokenTypeLoginSession,
				}); err == nil {
					if len(user.([]*models.User)) >= 0 {
						req.Request = req.Request.WithContext(context.WithValue(req.Request.Context(), global.MetaUser, user))
						filterChan.ProcessFilter(req, resp)
						return
					}
				}
			} else {
				errorHandler(req.Request.Context(), errors.NewServerError(http.StatusBadRequest, "unknown authorization method"), resp)
			}
			return
		} else {
			query := req.Request.URL.Query()
			if query.Get("authKey") != "" {
				if err = httputil.UnmarshalURLValues(query, &authReq); err != nil {
					errorHandler(req.Request.Context(), errors.NewServerError(http.StatusBadRequest, "unknown exception"), resp)
					return
				}
			}
		}
		if len(authReq.Data.AuthKey) > 0 || len(authReq.Data.AuthSecret) > 0 {
			if authReq.Data.AuthSign != "" {
				if authReq.Data.Payload, err = sign.GetPayloadFromHTTPRequest(req.Request); err != nil {
					errorHandler(req.Request.Context(), errors.NewServerError(http.StatusBadRequest, "Failed to get payload"), resp)
				}
			}
			if user, err := endpoints.Authentication(req.Request.Context(), authReq); err == nil {
				if len(user.([]*models.User)) >= 0 {
					req.Request = req.Request.WithContext(context.WithValue(req.Request.Context(), global.MetaUser, user))
					filterChan.ProcessFilter(req, resp)
					return
				}
			}
		}

		if autoRedirectToLoginPage, ok := req.SelectedRoute().Metadata()[global.MetaAutoRedirectToLoginPage].(bool); ok && autoRedirectToLoginPage {
			resp.Header().Set("Location", fmt.Sprintf("/admin/user/login?redirect_uri=%s", url.QueryEscape(req.Request.RequestURI)))
			resp.WriteHeader(302)
			return
		}
		errorHandler(req.Request.Context(), errors.NewServerError(http.StatusUnauthorized, "Not logged in or identity expired"), resp)
		return
	}
}

func HTTPLoggingFilter(req *restful.Request, resp *restful.Response, filterChan *restful.FilterChain) {
	ctx := req.Request.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	traceId := req.HeaderParameter("TraceId")
	if len(traceId) > 36 || len(traceId) <= 0 {
		if traceId = req.HeaderParameter("X-Request-Id"); len(traceId) > 36 || len(traceId) <= 0 {
			traceId = logs.NewTraceId()
		}
	}
	logger := log.With(logs.GetRootLogger(), global.TraceIdName, traceId)
	req.Request = req.Request.WithContext(context.WithValue(context.WithValue(ctx, global.TraceIdName, traceId), global.LoggerName, logger))
	start := time.Now()

	defer func() {
		if r := recover(); r != nil {
			errorEncoder(req.Request.Context(), errors.NewServerError(http.StatusForbidden, "Server exception"), resp)
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
		logger = log.With(logger,
			"msg", "HTTP response send.",
			"title", "response",
			"[httpURI]", req.Request.RequestURI,
			"[status]", resp.StatusCode(),
			"[contentType]", resp.Header().Get("Content-Type"),
			"[contentLength]", resp.ContentLength(),
		)
		level.Info(logger).Log("[totalTime]", fmt.Sprintf("%dms", time.Since(start).Milliseconds()))
	}()
	preBuf := w.Must[buffer.PreReader](buffer.NewPreReader(req.Request.Body, 1024))
	req.Request.Body = preBuf
	level.Info(logger).Log(
		"msg", "HTTP request received.",
		"title", "request",
		"[httpURI]", fmt.Sprintf("%s %s %s", req.Request.Method, req.Request.RequestURI, req.Request.Proto),
		"[contentType]", req.HeaderParameter("Content-Type"),
		"[header]", httputil.Map[string, []string](req.Request.Header),
		"[contentLength]", req.Request.ContentLength,
		"[body]", preBuf,
	)
	filterChan.ProcessFilter(req, resp)
}
