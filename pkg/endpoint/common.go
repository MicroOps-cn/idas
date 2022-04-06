package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"idas/pkg/service"
	"mime/multipart"
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
		var (
			f       multipart.File
			fileKey string
			data    = make(map[string]string)
		)
		for fileName, fhs := range stdReq.MultipartForm.File {
			if len(fhs) > 0 {
				f, err = fhs[0].Open()
				if err != nil {
					return nil, err
				} else if fileKey, err = s.UploadFile(fileName, f); err != nil {
					return nil, err
				}
				data[fileName] = fileKey
			}
		}
		resp.Data = data
		//resp.Data, resp.Total, resp.Error = s.GetUsers(ctx)
		return &resp, nil
	}
}
