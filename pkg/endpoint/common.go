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

package endpoint

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/MicroOps-cn/fuck/log"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log/level"
	uuid "github.com/satori/go.uuid"

	"github.com/MicroOps-cn/idas/pkg/service"
)

type FileUploadRequest struct{}

func MakeUploadFileEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp := SimpleResponseWrapper[interface{}]{}
		stdReq := request.(RestfulRequester).GetRestfulRequest().Request
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
		return &resp, nil
	}
}

type FileDownloadRequest struct {
	Id       string `json:"id"`
	Download bool   `json:"download"`
}

func MakeDownloadFileEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		logger := log.GetContextLogger(ctx)
		req := request.(Requester).GetRequestData().(*FileDownloadRequest)
		stdResp := request.(RestfulRequester).GetRestfulResponse()
		var (
			f        io.ReadCloser
			mimiType string
			fileName string
		)
		f, mimiType, fileName, err = s.DownloadFile(ctx, req.Id)
		if f != nil {
			defer f.Close()
		}
		if err != nil {
			return nil, err
		}
		if req.Download {
			stdResp.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
			stdResp.Header().Add("Content-Type", "application/octet-stream")
		} else if len(mimiType) != 0 {
			stdResp.Header().Add("Cache-Control", "max-age=86400")
			stdResp.Header().Add("Content-Type", mimiType)
		}
		_, err = io.Copy(stdResp, f)
		if err != nil {
			level.Error(logger).Log("msg", "failed to write response", "err", err)
		}
		return nil, nil
	}
}

func GetEventMeta(ctx context.Context, action string, beginTime time.Time, err error, resp interface{}) (eventId, message string, status bool, took time.Duration) {
	eventId = log.GetTraceId(ctx)
	if u, e := uuid.FromString(eventId); e == nil {
		eventId = u.String()
	}
	if err != nil {
		message = fmt.Sprintf("Calling the %s method failed, err: %s", action, err)
	} else if r, ok := resp.(endpoint.Failer); ok && r.Failed() != nil {
		message = fmt.Sprintf("Calling the %s method failed, err: %s", action, r.Failed())
	} else {
		message = fmt.Sprintf("Successfully called %s method.", action)
		status = true
	}
	return eventId, message, status, time.Since(beginTime)
}
