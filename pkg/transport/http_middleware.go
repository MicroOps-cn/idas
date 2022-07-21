package transport

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"idas/config"
	"idas/pkg/service/models"
	"idas/pkg/utils/httputil"
	"idas/pkg/utils/sets"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"idas/pkg/endpoint"
	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/logs"
)

var innerKey = sets.New[string]("authKey", "authSign", "authSecret", "authMethod", "authAlgorithm")

func GetPayload(r *http.Request) (string, error) {
	requestTime, err := time.Parse(r.Header.Get("date"), time.RFC1123)
	if err != nil {
		return "", err
	} else if time.Since(requestTime) > time.Minute*10 {
		return "", fmt.Errorf("request has expired")
	}
	var bodyHash string
	if r.ContentLength > 0 {
		contentType, _, _ := strings.Cut(r.Header.Get("content-type"), ";")
		if len(contentType) > 0 {
			switch contentType {
			case restful.MIME_JSON, restful.MIME_XML, "application/x-www-form-urlencoded":
				if r.ContentLength > config.Get().Global.MaxBodySize.Capacity {
					body, err := ioutil.ReadAll(r.Body)
					r.Body.Close()
					if err != nil {
						return "", err
					} else {
						bodyHash = fmt.Sprintf("%x", md5.Sum(body))
						r.Body = io.NopCloser(bytes.NewBuffer(body))
					}
				}
			}
		}
	}
	if len(bodyHash) == 0 {
		bodyHash = r.Header.Get("x-body-hash")
	}
	payload := strings.Builder{}
	payload.WriteString(r.Method + "\n")
	payload.WriteString(bodyHash + "\n")
	payload.WriteString(r.Header.Get("content-type") + "\n")
	payload.WriteString(r.Header.Get("date") + "\n")
	var urlQuery = url.Values{}
	for key, value := range r.URL.Query() {
		if !innerKey.Has(key) {
			for _, v := range value {
				urlQuery.Add(key, v)
			}
		}
	}
	payload.WriteString(r.URL.RawPath + "?" + urlQuery.Encode())
	return payload.String(), nil
}

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
			req.Request = req.Request.WithContext(context.WithValue(req.Request.Context(), global.LoginSession, strings.Split(loginSessionID.Value, ",")))
			if user, err := endpoints.GetLoginSession(req.Request.Context(), loginSessionID.Value); err == nil {
				if len(user.([]*models.User)) >= 0 {
					req.Request = req.Request.WithContext(context.WithValue(req.Request.Context(), global.MetaUser, user))
					filterChan.ProcessFilter(req, resp)
					return
				}
			}
		}

		authReq := HttpRequest[endpoint.AuthenticationRequest]{}
		if username, password, ok := req.Request.BasicAuth(); ok {
			authReq.Data.AuthKey = username
			authReq.Data.AuthSecret = password
		} else if auth := req.Request.Header.Get("Authorization"); len(auth) != 0 {
			errorHandler(req.Request.Context(), errors.NewServerError(http.StatusBadRequest, "unknown authorization method"), resp)
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
				if authReq.Data.Payload, err = GetPayload(req.Request); err != nil {
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
	level.Info(logger).Log(
		"msg", "HTTP request received.",
		"title", "request",
		"[httpURI]", req.Request.RequestURI,
		"[method]", req.Request.Method,
		"[proto]", req.Request.Proto,
		"[contentType]", req.HeaderParameter("Content-Type"),
		"[contentLength]", req.Request.ContentLength,
	)

	defer func() {
		if r := recover(); r != nil {
			errorEncoder(req.Request.Context(), errors.NewServerError(http.StatusForbidden, "Server exception"), resp)
			buffer := bytes.NewBufferString(fmt.Sprintf("recover from panic situation: - %v\n", r))
			for i := 2; ; i++ {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				buffer.WriteString(fmt.Sprintf("    %s:%d\n", file, line))
			}
			level.Error(logger).Log("msg", buffer.String())
		}
		logger = log.With(logger,
			"msg", "HTTP response send.",
			"title", "response",
			"[httpURI]", req.Request.RequestURI,
			"[status]", resp.StatusCode(),
			"[contentType]", resp.Header().Get("Content-Type"),
			"[contentLength]", resp.ContentLength(),
		)
		level.Info(logger).Log("[totalTime]", time.Since(start)/time.Millisecond)
	}()
	filterChan.ProcessFilter(req, resp)
}
