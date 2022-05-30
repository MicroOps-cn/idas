package transport

import (
	"bytes"
	"context"
	"fmt"
	"idas/pkg/service/models"
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

func HTTPLoginAuthentication(endpoints endpoint.Set) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, filterChan *restful.FilterChain) {
		if req.SelectedRoute() == nil {
			errorEncoder(req.Request.Context(), errors.NewServerError(http.StatusNotFound, "Not Found: "+req.Request.RequestURI), resp)
			return
		}
		if needLogin, ok := req.SelectedRoute().Metadata()[global.MetaNeedLogin].(bool); ok && !needLogin {
			filterChan.ProcessFilter(req, resp)
			return
		}

		errorHandler := errorEncoder

		loginSessionID, err := req.Request.Cookie(global.LoginSession)
		if err == nil {
			req.SetAttribute(global.LoginSession, strings.Split(loginSessionID.Value, ","))
			if user, err := endpoints.GetLoginSession(req.Request.Context(), strings.Split(loginSessionID.Value, ",")); err == nil {
				if len(user.([]*models.User)) >= 0 {
					req.SetAttribute(global.AttrUser, user)
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
		errorHandler(req.Request.Context(), errors.NewServerError(http.StatusForbidden, "Not logged in or identity expired"), resp)
		return
	}
}

func HTTPLogging(req *restful.Request, resp *restful.Response, filterChan *restful.FilterChain) {
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
