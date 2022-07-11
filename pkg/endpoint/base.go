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

type Requester interface {
	GetRequestData() interface{}
}

type RestfulRequester interface {
	Requester
	GetRestfulRequest() *restful.Request
	GetRestfulResponse() *restful.Response
}

type BaseResponse[T any] struct {
	Error        error  `json:"-"`
	Data         T      `json:"data,omitempty"`
	ErrorMessage string `json:"errorMessage"`
}

func (l BaseResponse[T]) GetData() interface{} {
	return l.Data
}

func (l BaseResponse[T]) Failed() error {
	if len(l.ErrorMessage) != 0 {
		return errors.NewServerError(http.StatusOK, l.ErrorMessage)
	}
	return l.Error
}

type BaseListResponse[T any] struct {
	BaseTotalResponse[T]
	Current  int64 `json:"current,omitempty"`
	PageSize int64 `json:"pageSize,omitempty"`
}

func (b BaseListResponse[T]) GetPageSize() int64 {
	return b.PageSize
}

func (b BaseListResponse[T]) GetCurrent() int64 {
	return b.Current
}

func NewBaseListResponse[T any](req *BaseListRequest) BaseListResponse[T] {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.Current == 0 {
		req.Current = 1
	}
	return BaseListResponse[T]{
		PageSize: req.PageSize,
		Current:  req.Current,
		BaseTotalResponse: BaseTotalResponse[T]{
			BaseResponse: BaseResponse[T]{},
		},
	}
}

type BaseTotalResponse[T any] struct {
	BaseResponse[T]
	Total int64 `json:"total,omitempty"`
}

func (b BaseTotalResponse[T]) GetTotal() int64 {
	return b.Total
}
