package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/emicklei/go-restful/v3"
	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"

	"idas/pkg/endpoint"
	"idas/pkg/errors"
	"idas/pkg/global"
	"idas/pkg/logs"
	"idas/pkg/utils/buffer"
	"idas/pkg/utils/httputil"
)

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(endpoints endpoint.Set, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	if zipkinTracer != nil {
		// Zipkin HTTP Server Trace can either be instantiated per endpoint with a
		// provided operation name or a global tracing service can be instantiated
		// without an operation name and fed to each Go kit endpoint as ServerOption.
		// In the latter case, the operation name will be the endpoint's http method.
		// We demonstrate a global tracing service here.
		options = append(options, zipkin.HTTPServerTrace(zipkinTracer))
	}

	m := restful.NewContainer()
	options = append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Concat", logger)))
	InstallHTTPApi(logger, m, options, endpoints)
	return m
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	traceId := ctx.Value(global.TraceIdName).(string)
	resp := responseWrapper{
		ErrorMessage: err.Error(),
		TraceId:      traceId,
		Success:      false,
	}
	if serverErr, ok := err.(errors.ServerError); ok {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(serverErr.StatusCode())
		resp.ErrorCode = serverErr.Code()
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		level.Info(ctx.Value(global.LoggerName).(log.Logger)).Log("msg", "failed to write response")
	}
}

type responseWrapper struct {
	Success      bool        `json:"success"`
	Data         interface{} `json:"data,inline"`
	ErrorCode    string      `json:"errorCode,omitempty"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
	TraceId      string      `json:"traceId"`
	Current      int64       `json:"current,omitempty"`
	PageSize     int64       `json:"pageSize,omitempty"`
	Total        int64       `json:"total"`
}

// decodeHTTPRequest Decode HTTP requests into request types
func decodeHTTPRequest[RequestType any](_ context.Context, r *http.Request) (interface{}, error) {
	var req RequestType
	var err error
	restfulReq := r.Context().Value(global.RestfulRequestContextName).(*restful.Request)
	restfulResp := r.Context().Value(global.RestfulResponseContextName).(*restful.Response)
	logger := logs.GetContextLogger(r.Context())
	query := restfulReq.Request.URL.Query()
	if len(query) > 0 {
		if err = httputil.UnmarshalURLValues(query, &req); err != nil {
			return nil, fmt.Errorf("failed to decode url query：%s", err)
		}
	}
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		if err = r.ParseForm(); err != nil {
			return nil, fmt.Errorf("failed to parse form data：%s", err)
		} else if len(r.Form) > 0 {
			if err = httputil.UnmarshalURLValues(r.Form, &req); err != nil {
				return nil, fmt.Errorf("failed to decode form data：%s", err)
			}
		}
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		if err = r.ParseMultipartForm(1e6); err != nil {
			return nil, fmt.Errorf("failed to parse multipart form data：%s", err)
		} else if len(r.Form) > 0 {
			if err = httputil.UnmarshalURLValues(r.Form, &req); err != nil {
				return nil, fmt.Errorf("failed to decode multipart form data：%s", err)
			}
		}
	} else if len(contentType) > 0 {
		logWriter := logs.NewWriterAdapter(level.Debug(log.With(logger, "caller", log.Caller(6))), logs.Prefix("decode http request: ", true))
		if err = json.NewDecoder(io.TeeReader(r.Body, buffer.LimitWriter(logWriter, 1024, buffer.LimitWriterIgnoreError))).Decode(&req); err != nil {
			return nil, fmt.Errorf("failed to decode request body：%s", err)
		}
	}
	if len(restfulReq.PathParameters()) > 0 {
		if err = httputil.UnmarshalURLValues(httputil.MapToURLValues(restfulReq.PathParameters()), &req); err != nil {
			return nil, fmt.Errorf("failed to decode path parameters：%s", err)
		}
	}

	if rr, ok := interface{}(&req).(endpoint.RestfulRequester); ok {
		rr.SetRestfulRequest(restfulReq)
		rr.SetRestfulResponse(restfulResp)
	}
	level.Debug(logger).Log("msg", "decoded http request", "req", fmt.Sprintf("%#v", req))
	if ok, err := govalidator.ValidateStruct(req); err != nil {
		return &req, errors.NewServerError(http.StatusBadRequest, err.Error())
	} else if !ok {
		return &req, errors.NewServerError(http.StatusBadRequest, "params error")
	}
	return &req, err
}

// encodeHTTPResponse Encode the response as an HTTP response message
func encodeHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	logger := logs.GetContextLogger(ctx)
	if f, ok := response.(kitendpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	traceId := ctx.Value(global.TraceIdName).(string)
	resp := responseWrapper{Success: true, TraceId: traceId}
	if l, ok := response.(endpoint.Lister); ok {
		resp.Data = l.GetData()
		resp.Total = l.GetTotal()
		resp.PageSize = l.GetPageSize()
		resp.Current = l.GetCurrent()
	} else {
		if t, ok := response.(endpoint.Total); ok {
			resp.Total = t.GetTotal()
		}
		resp.Data = response
	}

	logWriter := logs.NewWriterAdapter(level.Debug(log.With(logger, "resp", fmt.Sprintf("%#v", resp), "caller", log.Caller(6))), logs.Prefix("encoded http response: ", true))
	return json.NewEncoder(io.MultiWriter(w, buffer.LimitWriter(logWriter, 1024, buffer.LimitWriterIgnoreError))).Encode(resp)
}
