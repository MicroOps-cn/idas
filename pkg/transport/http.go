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
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	stdlog "log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	logs "github.com/MicroOps-cn/fuck/log"
	"github.com/asaskevich/govalidator"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-openapi/spec"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/MicroOps-cn/idas/config"
	"github.com/MicroOps-cn/idas/pkg/endpoint"
	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/utils/buffer"
	"github.com/MicroOps-cn/idas/pkg/utils/httputil"
	w "github.com/MicroOps-cn/idas/pkg/utils/wrapper"
)

//go:embed static
var staticFs embed.FS

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(ctx context.Context, logger log.Logger, endpoints endpoint.Set, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, openapiPath string) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log.With(logger, "caller", logs.Caller(9)))),
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
	restful.TraceLogger(stdlog.New(log.NewStdlibAdapter(level.Info(logger)), "[restful]", stdlog.LstdFlags|stdlog.Lshortfile))
	m.Filter(HTTPLoggingFilter)
	var specTags []spec.Tag
	for _, serviceGenerator := range apiServiceSet {
		specTag, svcs := serviceGenerator(options, endpoints)
		for _, svc := range svcs {
			m.Add(svc)
		}
		specTags = append(specTags, specTag)
	}
	if openapiPath != "" {
		level.Info(logger).Log("msg", fmt.Sprintf("enable openapi on `%s`", openapiPath))
		specConf := restfulspec.Config{
			WebServices: m.RegisteredWebServices(),
			APIPath:     openapiPath,
			PostBuildSwaggerObjectHandler: func(swo *spec.Swagger) {
				swo.Info = &spec.Info{
					InfoProps: spec.InfoProps{
						Title:       "ItemTestService",
						Description: "Resource for managing ItemTests",
						Version:     "1.0.0",
					},
				}
				swo.Tags = specTags
			},
		}
		m.Add(restfulspec.NewOpenAPIService(specConf))
	}
	webPrefix := ctx.Value(global.HTTPWebPrefixKey).(string)
	m.Handle(webPrefix, http.StripPrefix(webPrefix, http.FileServer(http.FS(w.M[fs.FS](fs.Sub(staticFs, "static"))))))
	return m
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	logger := ctx.Value(global.LoggerName).(log.Logger)
	level.Error(logger).Log("err", err, "msg", "failed to http request")
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
		level.Info(logger).Log("msg", "failed to write response")
	}
}

type ResponseWrapper[T any] struct {
	Data T `json:"data"`
}

type responseWrapper struct {
	Success      bool        `json:"success"`
	Data         interface{} `json:"data"`
	ErrorCode    string      `json:"errorCode,omitempty"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
	TraceId      string      `json:"traceId"`
	Current      int64       `json:"current,omitempty"`
	PageSize     int64       `json:"pageSize,omitempty"`
	Total        int64       `json:"total"`
}

type HTTPRequest[T any] struct {
	Data            T `json:"data"`
	restfulRequest  *restful.Request
	restfulResponse *restful.Response
}

func (b HTTPRequest[T]) GetRequestData() interface{} {
	return &b.Data
}

func (b HTTPRequest[T]) GetRestfulRequest() *restful.Request {
	return b.restfulRequest
}

func (b HTTPRequest[T]) GetRestfulResponse() *restful.Response {
	return b.restfulResponse
}

var _ endpoint.RestfulRequester = &HTTPRequest[any]{}

func isProtoMessage(v interface{}) (proto.Message, bool) {
	msg, ok := v.(proto.Message)
	return msg, ok
}

func valid(data interface{}) (bool, error) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Struct:
		return govalidator.ValidateStruct(data)
	case reflect.Slice:
		valOf := reflect.ValueOf(data)
		for i := 0; i < valOf.Len(); i++ {
			b, err := valid(valOf.Index(i).Interface())
			if err != nil || !b {
				return b, err
			}
		}
	}
	return true, nil
}

// decodeHTTPRequest Decode HTTP requests into request types
func decodeHTTPRequest[RequestType any](_ context.Context, stdReq *http.Request) (interface{}, error) {
	restfulReq := stdReq.Context().Value(global.RestfulRequestContextName).(*restful.Request)
	restfulResp := stdReq.Context().Value(global.RestfulResponseContextName).(*restful.Response)
	req := HTTPRequest[RequestType]{restfulRequest: restfulReq, restfulResponse: restfulResp}
	var err error
	logger := logs.GetContextLogger(stdReq.Context())
	r := restfulReq.Request
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE" {
		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "multipart/form-data") {
			restfulReq.Request.Body = http.MaxBytesReader(restfulResp.ResponseWriter, r.Body, config.Get().Global.MaxUploadSize.Capacity)
			if err = restfulReq.Request.ParseMultipartForm(config.Get().Global.MaxBodySize.Capacity); err != nil {
				return nil, errors.NewServerError(http.StatusBadRequest, "request too large")
			}
		} else {
			r.Body = http.MaxBytesReader(restfulResp.ResponseWriter, r.Body, config.Get().Global.MaxBodySize.Capacity)
			if contentType == "application/x-www-form-urlencoded" {
				if err = r.ParseForm(); err != nil {
					return nil, fmt.Errorf("failed to parse form data：%s", err)
				} else if len(r.Form) > 0 {
					if err = httputil.UnmarshalURLValues(r.Form, &req.Data); err != nil {
						return nil, fmt.Errorf("failed to decode form data: data=%s, err=%s", r.Form, err)
					}
				}
			} else if contentType == restful.MIME_JSON {
				if data, ok := isProtoMessage(&req.Data); ok {
					logWriter := logs.NewWriterAdapter(level.Debug(log.With(logger, "caller", logs.Caller(12))), logs.Prefix("decode http request: ", true))
					if err = jsonpb.Unmarshal(io.TeeReader(r.Body, buffer.LimitWriter(logWriter, 1024, buffer.LimitWriterIgnoreError)), data); err != nil {
						return nil, fmt.Errorf("failed to decode request body：%s", err)
					}
				} else {
					logWriter := logs.NewWriterAdapter(level.Debug(log.With(logger, "caller", logs.Caller(9))), logs.Prefix("decode http request: ", true))
					if err = json.NewDecoder(io.TeeReader(r.Body, buffer.LimitWriter(logWriter, 1024, buffer.LimitWriterIgnoreError))).Decode(&req.Data); err != nil {
						return nil, fmt.Errorf("failed to decode request body：%s", err)
					}
				}
			}
		}
	}

	query := restfulReq.Request.URL.Query()
	if len(query) > 0 {
		if err = httputil.UnmarshalURLValues(query, &req.Data); err != nil {
			return nil, fmt.Errorf("failed to decode url query: %s", err)
		}
	}
	if len(restfulReq.PathParameters()) > 0 {
		if err = httputil.UnmarshalURLValues(httputil.MapToURLValues(restfulReq.PathParameters()), &req); err != nil {
			return nil, fmt.Errorf("failed to decode path parameters：%s", err)
		}
	}

	req.restfulRequest = restfulReq
	req.restfulResponse = restfulResp
	level.Debug(logger).Log("msg", "decoded http request", "req", string(w.Must[[]byte](json.Marshal(req))))
	if ok, err := valid(req.Data); err != nil {
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
	} else if response != nil {
		if t, ok := response.(endpoint.Total); ok {
			resp.Total = t.GetTotal()
		}
		if t, ok := response.(endpoint.HasData); ok {
			resp.Data = t.GetData()
		} else {
			resp.Data = response
		}
	}

	logWriter := logs.NewWriterAdapter(level.Debug(log.With(logger, "resp", fmt.Sprintf("%#v", resp), "caller", logs.Caller(7))), logs.Prefix("encoded http response: ", true))
	return json.NewEncoder(io.MultiWriter(w, buffer.LimitWriter(logWriter, 1024, buffer.LimitWriterIgnoreError))).Encode(resp)
}

func simpleEncodeHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	logger := logs.GetContextLogger(ctx)
	if f, ok := response.(kitendpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	logWriter := logs.NewWriterAdapter(level.Debug(log.With(logger, "resp", fmt.Sprintf("%#v", response), "caller", logs.Caller(7))), logs.Prefix("encoded http response: ", true))
	return json.NewEncoder(io.MultiWriter(w, buffer.LimitWriter(logWriter, 1024, buffer.LimitWriterIgnoreError))).Encode(response)
}

func WrapHTTPHandler(h *httptransport.Server) func(*restful.Request, *restful.Response) {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := req.Request.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		request := req.Request.WithContext(context.WithValue(context.WithValue(ctx, global.RestfulResponseContextName, resp), global.RestfulRequestContextName, req))
		h.ServeHTTP(resp, request)
	}
}

func NewKitHTTPServer[RequestType any](dp kitendpoint.Endpoint, options []httptransport.ServerOption) restful.RouteFunction {
	return WrapHTTPHandler(httptransport.NewServer(
		dp,
		decodeHTTPRequest[RequestType],
		encodeHTTPResponse,
		options...,
	))
}

func NewSimpleKitHTTPServer[RequestType any](
	dp kitendpoint.Endpoint,
	dec httptransport.DecodeRequestFunc,
	enc httptransport.EncodeResponseFunc, options []httptransport.ServerOption,
) restful.RouteFunction {
	return WrapHTTPHandler(httptransport.NewServer(
		dp,
		dec,
		enc,
		options...,
	))
}

const QueryTypeKey = "__query_type__"

func NewWebService(rootPath string, gv schema.GroupVersion, doc string) *restful.WebService {
	webservice := restful.WebService{}
	webservice.Path(rootPath + "/" + gv.Version + "/" + gv.Group).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).Doc(doc)
	return &webservice
}

func NewSimpleWebService(rootPath string, doc string) *restful.WebService {
	webservice := restful.WebService{}
	webservice.Path(rootPath).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).Doc(doc)
	return &webservice
}

const rootPath = "/api"

func StructToQueryParams(obj interface{}, nameFilter ...string) []*restful.Parameter {
	var params []*restful.Parameter
	typeOfObj := reflect.TypeOf(obj)
	valueOfObj := reflect.ValueOf(obj)
	// 通过 #NumField 获取结构体字段的数量
loopObjFields:
	for i := 0; i < typeOfObj.NumField(); i++ {
		field := typeOfObj.Field(i)

		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			params = append(params, StructToQueryParams(valueOfObj.Field(i).Interface(), nameFilter...)...)
		} else {
			if len(nameFilter) > 0 {
				for _, name := range nameFilter {
					if name == field.Name {
						goto handleField
					}
				}
				continue loopObjFields
			}
		handleField:
			jsonTag := strings.Split(field.Tag.Get("json"), ",")
			if len(jsonTag) > 0 && jsonTag[0] != "-" && jsonTag[0] != "" {
				param := restful.QueryParameter(
					jsonTag[0],
					field.Tag.Get("description"),
				).DataType(field.Type.String())
				if len(jsonTag) > 1 && jsonTag[1] == "omitempty" {
					param.Required(false)
				} else {
					param.Required(true)
				}
				if tag := field.Tag.Get("enum"); tag != "" {
					enums := map[string]string{}
					for idx, s := range strings.Split(tag, "|") {
						enums[strconv.Itoa(idx)] = s
					}
					param.AllowableValues(enums)
				} else if protoTag := field.Tag.Get("protobuf"); protoTag != "" {
					var typeName string
					for _, s := range strings.Split(protoTag, ",") {
						if strings.HasPrefix(s, "enum=") {
							typeName = s[5:]
							break
						}
					}
					if len(typeName) != 0 {
						enumMap := proto.EnumValueMap(typeName)
						enums := make(map[string]string, len(enumMap))
						for v, idx := range enumMap {
							enums[strconv.Itoa(int(idx))] = v
						}
						param.AllowableValues(enums)
						param.AddExtension("$ref", typeName)
						param.DataType("string")
						param.DataFormat("string")
					}
				}
				params = append(params, param)
			}
		}
	}
	return params
}

func NewProxyHandler(c context.Context, logger log.Logger, endpoints endpoint.Set, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) http.Handler {
	m := restful.NewContainer()
	m.Filter(HTTPLoggingFilter)
	m.Filter(HTTPProxyAuthenticationFilter(c, endpoints))
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log.With(logger, "caller", logs.Caller(9)))),
	}
	if zipkinTracer != nil {
		// Zipkin HTTP Server Trace can either be instantiated per endpoint with a
		// provided operation name or a global tracing service can be instantiated
		// without an operation name and fed to each Go kit endpoint as ServerOption.
		// In the latter case, the operation name will be the endpoint's http method.
		// We demonstrate a global tracing service here.
		options = append(options, zipkin.HTTPServerTrace(zipkinTracer))
	}
	options = append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Concat", logger)))
	m.HandleWithFilter("/", httptransport.NewServer(
		endpoints.ProxyRequest,
		func(_ context.Context, request *http.Request) (interface{}, error) {
			return request, nil
		},
		func(ctx context.Context, writer http.ResponseWriter, resp interface{}) error {
			if r, ok := resp.(*endpoint.ProxyResponse); ok {
				if r.Header != nil {
					for name, vals := range r.Header {
						for i, val := range vals {
							if i == 0 {
								r.Header.Set(name, val)
							} else {
								r.Header.Add(name, val)
							}
						}
					}
				}
				writer.WriteHeader(r.Code)
				if r.Body != nil {
					if _, err := io.Copy(writer, r.Body); err != nil {
						level.Error(logs.GetContextLogger(ctx)).Log("msg", "failed to copy response", "err", err)
					}
				}
			}
			return nil
		},
		options...,
	))
	return m
}
