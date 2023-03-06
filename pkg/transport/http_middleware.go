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
	"fmt"
	"github.com/MicroOps-cn/fuck/buffer"
	"github.com/MicroOps-cn/fuck/log"
	w "github.com/MicroOps-cn/fuck/wrapper"
	"github.com/MicroOps-cn/idas/pkg/logs"
	"github.com/emicklei/go-restful/v3"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/MicroOps-cn/idas/pkg/endpoint"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
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
				fmt.Printf("%#v", session)
				if s := session.(*models.ProxySession); s != nil {
					ctx = context.WithValue(ctx, global.MetaUser, s.User)
					ctx = context.WithValue(ctx, global.MetaProxyConfig, s.Proxy)
					req.Request = req.Request.WithContext(ctx)
					filterChan.ProcessFilter(req, resp)
					return
				}
			}
		}
		if req.Request.URL.Path == "/-/oauth" {
			if authReq := GetProxyOAuthState(req.Request); authReq != nil {
				authReq.restfulResponse = resp
				authReq.restfulRequest = req
				if ar, err := endpoints.OAuthTokens(ctx, authReq); err != nil {
					errorEncoder(ctx, err, resp)
				} else if oar, ok := ar.(*endpoint.OAuthTokenResponse); ok && len(oar.AccessToken) != 0 {
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
		if needLogin, ok := ctx.Value(global.MetaNeedLogin).(bool); ok {
			if !needLogin {
				filterChan.ProcessFilter(req, resp)
				return
			}
		}

		if token := getTokenByRequest(req.Request); token != nil {
			ctx = context.WithValue(ctx, global.LoginSession, token.Token)
			req.Request = req.Request.WithContext(ctx)
			if user, err := endpoints.GetSessionByToken(ctx, token); err == nil {
				if len(user.(models.Users)) >= 0 {
					ctx = context.WithValue(ctx, global.MetaUser, user)
					req.Request = req.Request.WithContext(ctx)
					filterChan.ProcessFilter(req, resp)
					return
				}
				authError = errors.NewMultipleServerError(http.StatusUnauthorized, "can't get user by token")
			} else {
				authError = errors.WithServerError(http.StatusUnauthorized, err, "failed to get session by token")
			}
		}
		if authReq, err := getAuthReqByRequest(req.Request); err != nil {
			errorEncoder(ctx, err, resp)
			return
		} else if authReq != nil {
			if user, err := endpoints.Authentication(ctx, authReq); err == nil {
				if user != nil && len(user.(models.Users)) >= 0 {
					req.Request = req.Request.WithContext(context.WithValue(ctx, global.MetaUser, user))
					filterChan.ProcessFilter(req, resp)
					return
				}
			} else if err != nil {
				authError = errors.WithServerError(http.StatusUnauthorized, err, "failed to get session by auth request")
			}
		}

		if autoRedirectToLoginPage, ok := ctx.Value(global.MetaAutoRedirectToLoginPage).(bool); ok && autoRedirectToLoginPage {
			if loginURL, ok := ctx.Value(global.HTTPLoginURLKey).(string); ok && len(loginURL) > 0 {
				resp.Header().Set("Location", fmt.Sprintf("%s?redirect_uri=%s", loginURL, url.QueryEscape(req.Request.RequestURI)))
			} else {
				resp.Header().Set("Location", fmt.Sprintf("/admin/user/login?redirect_uri=%s", url.QueryEscape(req.Request.RequestURI)))
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
	header.Del("Cookie")
	return w.NewStringer(func() string {
		for _, cookie := range cookies {
			cookieVal := cookie.Value
			if cookie.Name == global.LoginSession {
				cookieVal = fmt.Sprintf("[sha256]%x", sha256.Sum256([]byte(cookie.Value)))
			}
			header.Add("Cookie", fmt.Sprintf(cookie.Name, cookieVal))
		}
		return w.JSONStringer(header).String()
	})
}

func HTTPContextFilter(pctx context.Context) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		ctx := req.Request.Context()
		if ctx == nil {
			ctx = pctx
		}
		if req.SelectedRoute() != nil && req.SelectedRoute().Metadata() != nil {
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
		ctx := req.Request.Context()
		if ctx == nil {
			ctx = pctx
		}
		hasSensitiveData, _ := ctx.Value(global.MetaSensitiveData).(bool)
		traceId := req.HeaderParameter("TraceId")
		if len(traceId) > 36 || len(traceId) <= 0 {
			if traceId = req.HeaderParameter("X-Request-Id"); len(traceId) > 36 || len(traceId) <= 0 {
				traceId = log.NewTraceId()
			}
		}
		var logger kitlog.Logger
		ctx, logger = log.NewContextLogger(ctx, log.WithTraceId(traceId))
		req.Request = req.Request.WithContext(ctx)
		start := time.Now()

		defer func() {
			if r := recover(); r != nil {
				errorEncoder(ctx, errors.NewServerError(http.StatusForbidden, "Server exception"), resp)
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
		}()
		var reqBody fmt.Stringer
		if !hasSensitiveData {
			preBuf := w.M[buffer.PreReader](buffer.NewPreReader(req.Request.Body, 1024))
			req.Request.Body = preBuf
			reqBody = preBuf
		} else {
			reqBody = bytes.NewBufferString("<body>")
		}
		getSafeHeader(req.Request)
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
