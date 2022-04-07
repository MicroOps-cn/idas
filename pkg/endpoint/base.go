package endpoint

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"

	"idas/pkg/errors"
)

type Lister interface {
	GetPageSize() int64
	GetCurrent() int64
	GetTotal() int64
	GetData() interface{}
}
type HasData interface {
	GetData() interface{}
}
type Total interface {
	GetTotal() int64
}

type RestfulRequester interface {
	GetRestfulRequest() *restful.Request
	SetRestfulRequest(r *restful.Request)
	GetRestfulResponse() *restful.Response
	SetRestfulResponse(c *restful.Response)
}

var _ RestfulRequester = &BaseRequest{}

type BaseRequest struct {
	restfulRequest  *restful.Request
	restfulResponse *restful.Response
}

func (b BaseRequest) GetRestfulRequest() *restful.Request {
	return b.restfulRequest
}

func (b *BaseRequest) SetRestfulRequest(r *restful.Request) {
	b.restfulRequest = r
}

func (b BaseRequest) GetRestfulResponse() *restful.Response {
	return b.restfulResponse
}

func (b *BaseRequest) SetRestfulResponse(r *restful.Response) {
	b.restfulResponse = r
}

type BaseListRequest struct {
	BaseRequest
	PageSize int64  `json:"pageSize"`
	Current  int64  `json:"current"`
	Keywords string `json:"keywords"`
}
type BaseResponse struct {
	Error        error       `json:"-"`
	Data         interface{} `json:"data,omitempty"`
	ErrorMessage string      `json:"errorMessage"`
}

func (l BaseResponse) GetData() interface{} {
	return l.Data
}

func (l BaseResponse) Failed() error {
	if len(l.ErrorMessage) != 0 {
		return errors.NewServerError(http.StatusOK, l.ErrorMessage)
	}
	return l.Error
}

type BaseListResponse struct {
	BaseResponse
	Current  int64       `json:"current,omitempty"`
	PageSize int64       `json:"pageSize,omitempty"`
	Total    int64       `json:"total,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func (b BaseListResponse) GetPageSize() int64 {
	return b.PageSize
}

func (b BaseListResponse) GetCurrent() int64 {
	return b.Current
}

func (b BaseListResponse) GetTotal() int64 {
	return b.Total
}

func (b BaseListResponse) GetData() interface{} {
	return b.Data
}

func NewBaseListResponse(req BaseListRequest) BaseListResponse {
	return BaseListResponse{
		PageSize: req.PageSize,
		Current:  req.Current,
	}
}

type BaseTotalResponse struct {
	BaseResponse
	Total int64 `json:"total,omitempty"`
}

func (b BaseTotalResponse) GetTotal() int64 {
	return b.Total
}