package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"idas/pkg/service"
)

type FileUploadRequest struct {
	BaseRequest
}

type FileUploadResponse struct {
	BaseResponse
}

func MakeUploadFileEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*FileUploadRequest)
		resp := FileUploadResponse{}
		stdReq := req.GetRestfulRequest().Request

		for fileName, fhs := range stdReq.MultipartForm.File {
			if len(fhs) > 0 {
				f, err := fhs[0].Open()
				if err != nil {
					return nil, err
				}
				var _, _ = fileName, f
			}
		}
		//resp.Data, resp.Total, resp.Error = s.GetUsers(ctx)
		return &resp, nil
	}
}
