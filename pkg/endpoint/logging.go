/*
 Copyright © 2023 MicroOps-cn.

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
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"

	"github.com/MicroOps-cn/idas/pkg/errors"
	"github.com/MicroOps-cn/idas/pkg/global"
	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

func MakeGetEventsEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetEventsRequest)
		resp := NewBaseListResponse[[]*models.Event](&req.BaseListRequest)
		startTime, err := time.Parse(time.RFC3339Nano, req.StartTime)
		if err != nil {
			return nil, errors.NewServerError(http.StatusBadRequest, fmt.Sprintf("Parameter Error: failed to parse startTime: %s", err))
		}
		endTime, err := time.Parse(time.RFC3339Nano, req.EndTime)
		if err != nil {
			return nil, errors.NewServerError(http.StatusBadRequest, fmt.Sprintf("Parameter Error: failed to parse endTime: %s", err))
		}
		resp.Total, resp.Data, resp.Error = svc.GetEvents(ctx, map[string]string{}, req.Keywords, startTime, endTime, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeGetEventLogsEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetEventLogsRequest)
		resp := NewBaseListResponse[[]*models.EventLog](&req.BaseListRequest)
		resp.Total, resp.Data, resp.Error = svc.GetEventLogs(ctx, map[string]string{"event_id": req.EventId}, req.Keywords, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeGetCurrentUserEventsEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetCurrentUserEventsRequest)
		resp := NewBaseListResponse[[]*models.Event](&req.BaseListRequest)
		user, ok := ctx.Value(global.MetaUser).(*models.User)
		if !ok || user == nil {
			resp.Error = errors.NotLoginError()
			return &resp, nil
		}
		startTime, err := time.Parse(time.RFC3339Nano, req.StartTime)
		if err != nil {
			return nil, errors.NewServerError(http.StatusBadRequest, fmt.Sprintf("Parameter Error: failed to parse startTime: %s", err))
		}
		endTime, err := time.Parse(time.RFC3339Nano, req.EndTime)
		if err != nil {
			return nil, errors.NewServerError(http.StatusBadRequest, fmt.Sprintf("Parameter Error: failed to parse endTime: %s", err))
		}
		resp.Total, resp.Data, resp.Error = svc.GetEvents(ctx, map[string]string{"user_id": user.Id}, req.Keywords, startTime, endTime, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeGetCurrentUserEventLogsEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetCurrentUserEventLogsRequest)
		resp := NewBaseListResponse[[]*models.EventLog](&req.BaseListRequest)
		user, ok := ctx.Value(global.MetaUser).(*models.User)
		if !ok || user == nil {
			resp.Error = errors.NotLoginError()
			return &resp, nil
		}
		resp.Total, resp.Data, resp.Error = svc.GetEventLogs(ctx, map[string]string{"event_id": req.EventId, "user_id": user.Id}, req.Keywords, req.Current, req.PageSize)
		return &resp, nil
	}
}
