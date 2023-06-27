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
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/gogo/protobuf/proto"

	"github.com/MicroOps-cn/idas/pkg/service"
	"github.com/MicroOps-cn/idas/pkg/service/models"
)

func MakeDeletePageEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeletePageRequest)
		resp := SimpleResponseWrapper[interface{}]{}
		resp.Error = svc.DeletePages(ctx, []string{req.Id})
		return &resp, nil
	}
}

func MakeUpdatePageEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UpdatePageRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		page := &models.PageConfig{
			Model:       models.Model{Id: req.Id},
			Name:        req.Name,
			Description: req.Description,
			Fields:      req.Fields,
			Icon:        req.Icon,
		}
		resp.Error = svc.UpdatePage(ctx, page)
		return resp, nil
	}
}

func MakeCreatePageEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreatePageRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		page := &models.PageConfig{
			Name:        req.Name,
			Description: req.Description,
			Fields:      req.Fields,
			Icon:        req.Icon,
		}
		resp.Error = svc.CreatePage(ctx, page)
		return resp, nil
	}
}

func MakeGetPagesEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetPagesRequest)
		resp := NewBaseListResponse[[]*models.PageConfig](&req.BaseListRequest)
		filter := map[string]interface{}{}
		switch req.Status {
		case models.PageStatus_disabled:
			filter["is_disable"] = true
		case models.PageStatus_enabled:
			filter["is_disable"] = false
		}
		resp.Total, resp.Data, resp.BaseResponse.Error = svc.GetPages(ctx, filter, req.Keywords, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeGetPageEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetPageRequest)
		resp := SimpleResponseWrapper[*models.PageConfig]{}
		resp.Data, resp.Error = svc.GetPage(ctx, req.Id)
		return &resp, nil
	}
}

type PatchPagesRequest []PatchPageRequest

func (m *PatchPagesRequest) Reset()         { *m = PatchPagesRequest{} }
func (m *PatchPagesRequest) String() string { return proto.CompactTextString(m) }

func (m PatchPagesRequest) ProtoMessage() {}

func MakePatchPagesEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchPagesRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		var pages []map[string]interface{}
		for _, pageRequest := range *req {
			pagePatch := map[string]interface{}{
				"id": pageRequest.Id,
			}
			if pageRequest.Name != nil {
				pagePatch["name"] = pageRequest.IsDisable
			}
			if pageRequest.Description != nil {
				pagePatch["description"] = pageRequest.IsDisable
			}
			if pageRequest.Fields != nil {
				pagePatch["fields"] = pageRequest.IsDisable
			}
			if pageRequest.Icon != nil {
				pagePatch["icon"] = pageRequest.IsDisable
			}
			if pageRequest.IsDisable != nil {
				pagePatch["is_disable"] = pageRequest.IsDisable
			}
			if pageRequest.IsDelete != nil {
				pagePatch["delete_time"] = time.Now().UTC()
			}
			pages = append(pages, pagePatch)
		}
		resp.Error = svc.PatchPages(ctx, pages)
		return resp, nil
	}
}

func MakeDeletePageDataEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*DeletePageDataRequest)
		resp := SimpleResponseWrapper[interface{}]{}

		pageDatasPatch := []models.PageData{{
			Model:  models.Model{Id: req.Id, IsDelete: true},
			PageId: req.PageId,
		}}
		resp.Error = svc.PatchPageDatas(ctx, pageDatasPatch)
		return &resp, nil
	}
}

func MakeUpdatePageDataEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*UpdatePageDataRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		resp.Error = svc.UpdatePageData(ctx, req.PageId, req.Id, req.Data)
		return resp, nil
	}
}

func MakeCreatePageDataEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*CreatePageDataRequest)
		resp := SimpleResponseWrapper[struct{}]{}
		resp.Error = svc.CreatePageData(ctx, req.PageId, req.Data)
		return resp, nil
	}
}

func MakeGetPageDatasEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetPageDatasRequest)
		resp := NewBaseListResponse[[]*models.PageData](&req.BaseListRequest)
		if req.Filters == nil {
			req.Filters = map[string]string{}
		}
		req.Filters["page_id"] = req.PageId
		resp.Total, resp.Data, resp.BaseResponse.Error = svc.GetPageDatas(ctx, req.Filters, req.Keywords, req.Current, req.PageSize)
		return &resp, nil
	}
}

func MakeGetPageDataEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*GetPageDataRequest)
		resp := SimpleResponseWrapper[*models.PageData]{}
		resp.Data, resp.Error = svc.GetPageData(ctx, req.PageId, req.Id)
		return &resp, nil
	}
}

type PatchPageDatasRequest []PatchPageDataRequest

func (m *PatchPageDatasRequest) Reset()         { *m = PatchPageDatasRequest{} }
func (m *PatchPageDatasRequest) String() string { return proto.CompactTextString(m) }

func (m PatchPageDatasRequest) ProtoMessage() {}

func MakePatchPageDatasEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Requester).GetRequestData().(*PatchPageDatasRequest)
		pageId := request.(RestfulRequester).GetRestfulRequest().PathParameter("pageId")
		resp := SimpleResponseWrapper[struct{}]{}
		var pageDatasPatch []models.PageData
		for _, patch := range *req {
			pageDatasPatch = append(pageDatasPatch, models.PageData{
				Model:  models.Model{Id: patch.Id, IsDelete: patch.IsDelete},
				PageId: pageId,
				Data:   (*models.JSON)(patch.Data),
			})
		}
		resp.Error = svc.PatchPageDatas(ctx, pageDatasPatch)
		return resp, nil
	}
}
