package endpoint

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"idas/pkg/service"
	"io"
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
				} else if fileKey, err = s.UploadFile(ctx, fileName, fhs[0].Header.Get("Content-Type"), f); err != nil {
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

type FileDownloadRequest struct {
	BaseRequest
	Id       string `json:"id"`
	Download bool   `json:"download"`
}

func MakeDownloadFileEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*FileDownloadRequest)
		stdResp := req.GetRestfulResponse()
		var (
			f        io.ReadCloser
			mimiType string
			fileName string
		)
		f, mimiType, fileName, err = s.DownloadFile(ctx, req.Id)
		defer f.Close()
		if err != nil {
			return nil, err
		}
		if req.Download {
			stdResp.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
			stdResp.Header().Add("Content-Type", "application/octet-stream")
		} else if len(mimiType) != 0 {
			stdResp.Header().Add("Content-Type", mimiType)
		}
		io.Copy(stdResp, f)
		//resp.Data, resp.Total, resp.Error = s.GetUsers(ctx)
		return nil, nil
	}
}
