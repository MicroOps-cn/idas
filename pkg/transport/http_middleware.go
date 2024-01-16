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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/MicroOps-cn/fuck/buffer"
	http2 "github.com/MicroOps-cn/fuck/http"
	"github.com/MicroOps-cn/fuck/log"
	"github.com/MicroOps-cn/fuck/signals"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/MicroOps-cn/idas/config"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-kit/kit/metrics/prometheus"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/MicroOps-cn/idas/pkg/endpoint"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/logs"
	"github.com/MicroOps-cn/idas/pkg/service/models"
	"github.com/MicroOps-cn/idas/pkg/utils/httputil"
	"github.com/MicroOps-cn/idas/pkg/utils/sign"
)

var proxyRedirectTemplate *template.Template

func init() {
	var err error
	proxyRedirectTemplate, err = template.ParseFS(templFs, "template/proxy_redirect_to_login.html")
	if err != nil {
		fmt.Println("failed to parse template: proxy_redirect_to_login.html")
		os.Exit(1)
	}
}

func getTokenByRequest(req *http.Request) *endpoint.GetSessionParams {
	loginSessionID, err := req.Cookie(global.LoginSession)
	if err == nil {
		return &endpoint.GetSessionParams{
			Token:     loginSessionID.Value,
			TokenType: models.TokenTypeLoginSession,
		}
	}
	if auth := req.Header.Get("Authorization"); len(auth) != 0 {
		if strings.HasPrefix(auth, "Bearer ") {
			return &endpoint.GetSessionParams{
				Token:     strings.TrimPrefix(auth, "Bearer "),
				TokenType: models.TokenTypeToken,
			}
		}
	}
	return nil
}

func getAuthReqByRequest(req *http.Request) (*HTTPRequest[endpoint.AuthenticationRequest], error) {
	var err error
	authReq := &HTTPRequest[endpoint.AuthenticationRequest]{}
	if username, password, ok := req.BasicAuth(); ok {
		authReq.Data.AuthKey = username
		authReq.Data.AuthSecret = password
	} else {
		query := req.URL.Query()
		if query.Get("authKey") != "" {
			if err = httputil.UnmarshalURLValues(query, &authReq); err != nil {
				return nil, errors.NewServerError(400, "unknown exception")
			}
		}
	}
	if len(authReq.Data.AuthKey) > 0 || len(authReq.Data.AuthSecret) > 0 {
		if authReq.Data.AuthSign != "" {
			if authReq.Data.Payload, err = sign.GetPayloadFromHTTPRequest(req); err != nil {
				return nil, errors.ParameterError("Failed to get payload")
			}
		}
		return authReq, nil
	}
	return nil, nil
}

func GetProxyOAuthState(r *http.Request) *HTTPRequest[endpoint.OAuthTokenRequest] {
	logger := log.GetContextLogger(r.Context())
	var err error
	oauthState, err := r.Cookie(global.OAuthStateCookieName)
	if err != nil {
		level.Debug(logger).Log("msg", fmt.Sprintf("failed to get cookie <%s>", global.OAuthStateCookieName), "err", err)
		return nil
	}
	redirect, err := r.Cookie(global.RedirectURICookieName)
	if err != nil {
		level.Debug(logger).Log("msg", fmt.Sprintf("failed to get cookie <%s>", global.RedirectURICookieName), "err", err)
		return nil
	}

	redirectURI, err := url.QueryUnescape(redirect.Value)
	if err != nil {
		level.Debug(logger).Log("msg", fmt.Sprintf("failed to parse Unescape uri <%s>", redirect.Value), "err", err)
		return nil
	}
	clientId, err := r.Cookie(global.ClientIDCookieName)
	if err != nil {
		level.Debug(logger).Log("msg", fmt.Sprintf("failed to get cookie <%s>", global.ClientIDCookieName), "err", err)
		return nil
	}
	q := r.URL.Query()
	state := q.Get("state")
	if len(state) == 0 {
		level.Debug(logger).Log("msg", "state params is null")
		return nil
	}
	code := q.Get("code")
	if len(code) == 0 {
		level.Debug(logger).Log("msg", "code params is null")
		return nil
	}
	if hashStateCode(state, clientId.Value) != oauthState.Value {
		level.Debug(logger).Log("msg", "state is invalid", "clientId", clientId.Value, "state", state, "oauthState", oauthState.Value)
		return nil
	}

	return &HTTPRequest[endpoint.OAuthTokenRequest]{
		Data: endpoint.OAuthTokenRequest{
			State:               state,
			Code:                code,
			ClientId:            clientId.Value,
			RedirectUri:         redirectURI,
			DisableRefreshToken: true,
			GrantType:           endpoint.OAuthGrantType_proxy,
			TokenType:           endpoint.OAuthTokenType_Cookie,
		},
	}
}

func hashStateCode(state, clientId string) string {
	hashBytes := sha256.Sum256([]byte(state + state[1:5] + clientId))
	return hex.EncodeToString(hashBytes[:])
}

func HTTPProxyAuthenticationFilter(ptx context.Context, endpoints endpoint.Set) restful.FilterFunction {
	httpExternalURL, ok := ptx.Value(global.HTTPExternalURLKey).(string)
	if !ok {
		panic("system error")
	}
	return func(req *restful.Request, resp *restful.Response, filterChan *restful.FilterChain) {
		ctx := req.Request.Context()
		token := getTokenByRequest(req.Request)

		if req.Request.URL.Path == "/-/logout" {
			_, err := endpoints.UserLogout(ctx, HTTPRequest[any]{restfulRequest: req, restfulResponse: resp})
			if err != nil {
				errorEncoder(ctx, err, resp)
			}
			return
		}
		if token != nil {
			ctx = context.WithValue(ctx, global.LoginSession, token.Token)
			req.Request = req.Request.WithContext(ctx)
			if session, err := endpoints.GetProxySessionByToken(ctx, token); err == nil {
				if s := session.(*models.ProxySession); s != nil {
					ctx = context.WithValue(ctx, global.MetaUser, s.User)
					ctx = context.WithValue(ctx, global.MetaProxyConfig, s.Proxy)
					req.Request = req.Request.WithContext(ctx)
					filterChan.ProcessFilter(req, resp)
					return
				}
			} else if err == gobreaker.ErrOpenState || err == gobreaker.ErrTooManyRequests {
				errorEncoder(ctx, err, resp)
				return
			}
		}
		if req.Request.URL.Path == "/-/oauth" {
			if authReq := GetProxyOAuthState(req.Request); authReq != nil {
				authReq.restfulResponse = resp
				authReq.restfulRequest = req
				if ar, err := endpoints.OAuthTokens(ctx, authReq); err != nil {
					errorEncoder(ctx, err, resp)
				} else if oar, ok := ar.(*endpoint.OAuthTokenResponse); ok && len(oar.AccessToken) != 0 {
					for _, cookie := range oar.Cookies {
						resp.ResponseWriter.Header().Add("Set-Cookie", cookie)
					}
					for name, value := range oar.Headers {
						resp.ResponseWriter.Header().Add(name, value)
					}
					http.SetCookie(resp.ResponseWriter, &http.Cookie{Name: global.LoginSession, Path: "/", Value: oar.AccessToken})
					http.SetCookie(resp.ResponseWriter, &http.Cookie{Name: global.ClientIDCookieName, Path: "/", MaxAge: -1})
					http.SetCookie(resp.ResponseWriter, &http.Cookie{Name: global.OAuthStateCookieName, Path: "/", MaxAge: -1})
					http.SetCookie(resp.ResponseWriter, &http.Cookie{Name: global.RedirectURICookieName, Path: "/", MaxAge: -1})
					http.Redirect(resp.ResponseWriter, req.Request, authReq.Data.RedirectUri, http.StatusFound)
				}
			}
			return
		} else if strings.HasSuffix(req.Request.URL.Path, "/favicon.ico") {
			errorEncoder(ctx, errors.NotFoundError(), resp)
			return
		}

		proxyConfig, err := endpoints.GetProxyConfig(ctx, req)
		if err != nil {
			errorEncoder(ctx, err, resp)
			return
		}
		pc, ok := proxyConfig.(*endpoint.ProxyConfig)
		if !ok {
			errorEncoder(ctx, errors.NewServerError(http.StatusInternalServerError, "failed to parse proxy config"), resp)
			return
		}
		pc.ExternalURL = httpExternalURL
		pc.RedirectURICookieName = global.RedirectURICookieName
		tokenState := hashStateCode(pc.Token, pc.ClientId)
		http.SetCookie(resp.ResponseWriter, &http.Cookie{Name: global.OAuthStateCookieName, Path: "/", Value: tokenState})
		http.SetCookie(resp.ResponseWriter, &http.Cookie{Name: global.ClientIDCookieName, Path: "/", Value: pc.ClientId})
		err = proxyRedirectTemplate.Execute(resp.ResponseWriter, pc)
		if err != nil {
			level.Warn(log.GetContextLogger(ctx)).Log("err", err, "msg", "failed to response")
		}
	}
}

func HTTPApplicationAuthenticationFilter(endpoints endpoint.Set) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		ctx := req.Request.Context()
		var authReq *endpoint.AuthenticationRequest
		if username, password, ok := req.Request.BasicAuth(); ok {
			authReq = &endpoint.AuthenticationRequest{
				AuthKey:    username,
				AuthSecret: password,
			}
		} else {
			query := req.Request.URL.Query()
			clientId := query.Get("client_id")
			clientSecret := query.Get("client_secret")
			if len(clientId) != 0 && len(clientSecret) != 0 {
				authReq = &endpoint.AuthenticationRequest{
					AuthKey:    clientId,
					AuthSecret: clientSecret,
				}
			} else if req.Request.Method == "POST" {
				var buf bytes.Buffer
				var oauthReq endpoint.OAuthTokenRequest
				var err error
				contentType := httputil.GetContentType(req.Request.Header)
				if contentType == restful.MIME_JSON {
					req.Request.Body = io.NopCloser(io.TeeReader(req.Request.Body, &buf))
					err = json.NewDecoder(&buf).Decode(&oauthReq)
				} else if contentType == "application/x-www-form-urlencoded" {
					if err = req.Request.ParseForm(); err == nil {
						err = httputil.UnmarshalURLValues(req.Request.Form, &oauthReq)
					}
				}
				if (err == nil || err == io.EOF) && len(oauthReq.ClientId) != 0 && len(oauthReq.ClientSecret) != 0 {
					authReq = &endpoint.AuthenticationRequest{
						AuthKey:    oauthReq.ClientId,
						AuthSecret: oauthReq.ClientSecret,
					}
				}
			}
		}
		if authReq != nil {
			app, err := endpoints.AppAuthentication(ctx, &HTTPRequest[endpoint.AuthenticationRequest]{
				restfulRequest:  req,
				restfulResponse: resp,
				Data:            *authReq,
			})
			if err != nil {
				errorEncoder(ctx, err, resp)
			} else if app != nil {
				ctx = context.WithValue(ctx, global.MetaApp, app)
				req.Request = req.Request.WithContext(ctx)
				chain.ProcessFilter(req, resp)
				return
			}
		}
		errorEncoder(ctx, errors.NewServerError(http.StatusUnauthorized, "Not logged in or identity expired"), resp)
	}
}

func HTTPAuthenticationFilter(endpoints endpoint.Set) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, filterChan *restful.FilterChain) {
		ctx := req.Request.Context()
		var authError error
		if req.SelectedRoute() == nil {
			errorEncoder(ctx, errors.NewServerError(http.StatusNotFound, "Not Found: "+req.Request.RequestURI), resp)
			return
		}

		if token := getTokenByRequest(req.Request); token != nil {
			ctx = context.WithValue(ctx, global.LoginSession, token.Token)
			req.Request = req.Request.WithContext(ctx)
			sessionReq := &HTTPRequest[endpoint.GetSessionParams]{restfulRequest: req, restfulResponse: resp, Data: *token}
			if user, err := endpoints.GetSessionByToken(ctx, sessionReq); err == nil {
				if u := user.(*models.User); u != nil {
					ctx = context.WithValue(ctx, global.MetaUser, user)
					req.Request = req.Request.WithContext(ctx)
					filterChan.ProcessFilter(req, resp)
					return
				}
				authError = errors.NewServerError(http.StatusUnauthorized, "can't get user by token")
			} else {
				authError = errors.WithServerError(http.StatusUnauthorized, err, "failed to get session by token")
			}
		}
		if authReq, err := getAuthReqByRequest(req.Request); err != nil {
			errorEncoder(ctx, err, resp)
			return
		} else if authReq != nil {
			if user, err := endpoints.Authentication(ctx, authReq); err == nil {
				if user.(*models.User) != nil {
					req.Request = req.Request.WithContext(context.WithValue(ctx, global.MetaUser, user))
					filterChan.ProcessFilter(req, resp)
					return
				}
			} else if err != nil {
				authError = errors.WithServerError(http.StatusUnauthorized, err, "failed to get session by auth request")
			}
		}

		if needLogin, ok := ctx.Value(global.MetaNeedLogin).(bool); ok {
			if !needLogin {
				filterChan.ProcessFilter(req, resp)
				return
			}
		}
		if autoRedirectToLoginPage, ok := ctx.Value(global.MetaAutoRedirectToLoginPage).(bool); ok && autoRedirectToLoginPage {
			redirectURI := req.Request.RequestURI
			if externalURL, ok := ctx.Value(global.HTTPExternalURLKey).(string); ok {
				extURL, err := url.Parse(externalURL)
				if err == nil {
					extURL.Path = http2.JoinPath(extURL.Path, req.Request.URL.Path)
					extURL.RawQuery = req.Request.URL.RawQuery
					redirectURI = extURL.String()
				}
			}
			if loginURL, ok := ctx.Value(global.HTTPLoginURLKey).(string); ok && len(loginURL) > 0 {
				resp.Header().Set("Location", fmt.Sprintf("%s?redirect_uri=%s", loginURL, url.QueryEscape(redirectURI)))
			} else {
				resp.Header().Set("Location", fmt.Sprintf("/admin/account/login?redirect_uri=%s", url.QueryEscape(redirectURI)))
			}
			resp.WriteHeader(302)
			return
		}
		if authError != nil {
			errorEncoder(ctx, errors.WithServerError(http.StatusUnauthorized, authError, "Not logged in or identity expired"), resp)
		} else {
			errorEncoder(context.WithValue(ctx, DisableStackTrace, true), errors.NewServerError(http.StatusUnauthorized, "Not logged in or identity expired"), resp)
		}
	}
}

func getSafeHeader(req *http.Request) fmt.Stringer {
	header := req.Header.Clone()
	cookies := req.Cookies()
	return w.NewStringer(func() string {
		if auth := header.Get("Authorization"); len(auth) > 0 {
			header.Set("Authorization", fmt.Sprintf("[sha256]%x", sha256.Sum256([]byte(auth))))
		}
		header.Del("Cookie")
		for _, cookie := range cookies {
			cookieVal := cookie.Value
			if cookie.Name == global.LoginSession {
				cookieVal = fmt.Sprintf("[sha256]%x", sha256.Sum256([]byte(cookie.Value)))
			}
			header.Add("Cookie", fmt.Sprintf("%s=%s", cookie.Name, cookieVal))
		}
		return w.JSONStringer(header).String()
	})
}

var (
	requestsTotal = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Tracks the number of HTTP requests.",
	}, []string{"method", "code", "api"})
	requestDuration = prometheus.NewHistogramFrom(
		stdprometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Tracks the latencies for HTTP requests.",
			Buckets: stdprometheus.ExponentialBuckets(0.1, 3, 5),
		},
		[]string{"method", "code", "api"},
	)
)

func HTTPContextFilter(pctx context.Context) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		ch := signals.SignalHandler()
		ch.AddRequest(1)
		defer ch.DoneRequest()

		var route restful.RouteReader
		start := time.Now()
		defer func() {
			if route != nil {
				requestsTotal.With("method", route.Method(), "code", strconv.Itoa(resp.StatusCode()), "api", route.Path()).Add(1)
				requestDuration.With("method", route.Method(), "code", strconv.Itoa(resp.StatusCode()), "api", route.Path()).
					Observe(float64(time.Since(start) / time.Second))
			}
		}()
		ctx := req.Request.Context()
		if ctx == nil {
			ctx = pctx
		}
		if route = req.SelectedRoute(); route != nil && route.Metadata() != nil {
			metadata := req.SelectedRoute().Metadata()
			for key, val := range metadata {
				ctx = context.WithValue(ctx, key, val)
			}
		}
		req.Request = req.Request.WithContext(ctx)
		chain.ProcessFilter(req, resp)
	}
}

func HTTPLoggingFilter(pctx context.Context) func(req *restful.Request, resp *restful.Response, filterChan *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, filterChan *restful.FilterChain) {
		var logger kitlog.Logger
		ctx := req.Request.Context()
		if ctx == nil {
			ctx = pctx
		}
		hasSensitiveData, _ := ctx.Value(global.MetaSensitiveData).(bool)
		start := time.Now()
		spanName := req.Request.RequestURI
		var spanOptions []trace.SpanStartOption
		if req.SelectedRoute() != nil {
			spanName = req.SelectedRoute().Operation()
		}
		ctx, span := otel.GetTracerProvider().Tracer(config.Get().GetAppName()).Start(ctx, spanName, spanOptions...)
		traceId := span.SpanContext().TraceID()
		traceIdStr := traceId.String()
		if !traceId.IsValid() {
			traceIdStr = log.NewTraceId()
		}

		ctx, _ = log.NewContextLogger(ctx, log.WithTraceId(traceIdStr))
		req.Request = req.Request.WithContext(ctx)
		logger = log.GetContextLogger(ctx)
		defer func() {
			if r := recover(); r != nil {
				span.SetStatus(codes.Error, fmt.Sprintf("%+v", r))
				errorEncoder(ctx, errors.NewServerError(http.StatusInternalServerError, "Server exception"), resp)
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
				logs.WrapKeyName("httpURI"), req.Request.RequestURI,
				logs.WrapKeyName("status"), resp.StatusCode(),
				logs.WrapKeyName("contentType"), resp.Header().Get("Content-Type"),
				logs.WrapKeyName("contentLength"), resp.ContentLength(),
			)
			level.Info(logger).Log(logs.WrapKeyName("totalTime"), fmt.Sprintf("%dms", time.Since(start).Milliseconds()))
			span.End()
		}()
		var reqBody fmt.Stringer
		if !hasSensitiveData {
			preBuf := w.M[buffer.PreReader](buffer.NewPreReader(req.Request.Body, 1024))
			req.Request.Body = preBuf
			reqBody = preBuf
		} else {
			reqBody = bytes.NewBufferString("<body>")
		}
		level.Info(logger).Log(
			"msg", "HTTP request received.",
			logs.TitleKey, "request",
			logs.WrapKeyName("httpURI"), fmt.Sprintf("%s %s %s", req.Request.Method, req.Request.RequestURI, req.Request.Proto),
			logs.WrapKeyName("contentType"), req.HeaderParameter("Content-Type"),
			logs.WrapKeyName("header"), getSafeHeader(req.Request),
			logs.WrapKeyName("contentLength"), req.Request.ContentLength,
			logs.WrapKeyName("body"), reqBody,
		)
		filterChan.ProcessFilter(req, resp)
	}
}
