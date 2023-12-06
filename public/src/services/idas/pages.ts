// @ts-ignore

/* eslint-disable */
import { request } from '@/utils/request';

/** Get page list GET /api/v1/pages */
export async function getPages(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getPagesParams,
  options?: { [key: string]: any },
) {
  return request<API.GetPagesResponse>('/api/v1/pages', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Create page. POST /api/v1/pages */
export async function createPage(body: API.CreatePageRequest, options?: { [key: string]: any }) {
  return request<any>('/api/v1/pages', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Batch patch page config(Incremental). PATCH /api/v1/pages */
export async function patchPages(body: API.PatchPageRequest[], options?: { [key: string]: any }) {
  return request<API.BaseTotalResponse>('/api/v1/pages', {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Get a page configs. GET /api/v1/pages/${param0} */
export async function getPage(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getPageParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.GetPageResponse>(`/api/v1/pages/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** Update page (full). PUT /api/v1/pages/${param0} */
export async function updatePage(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.updatePageParams,
  body: API.UpdatePageRequest,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/pages/${param0}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** Delete a page. DELETE /api/v1/pages/${param0} */
export async function deletePage(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deletePageParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/pages/${param0}`, {
    method: 'DELETE',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** Get data list of page GET /api/v1/pages/${param0}/data */
export async function getPageDatas(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getPageDatasParams,
  options?: { [key: string]: any },
) {
  const { pageId: param0, ...queryParams } = params;
  return request<API.GetPageDatasResponse>(`/api/v1/pages/${param0}/data`, {
    method: 'GET',
    params: {
      ...queryParams,
    },
    ...(options || {}),
  });
}

/** Create a data of a page. POST /api/v1/pages/${param0}/data */
export async function createPageData(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.createPageDataParams,
  body: API.CreatePageDataRequest,
  options?: { [key: string]: any },
) {
  const { pageId: param0, ...queryParams } = params;
  return request<any>(`/api/v1/pages/${param0}/data`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** Batch patch data of a page(Incremental). PATCH /api/v1/pages/${param0}/data */
export async function patchPageDatas(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.patchPageDatasParams,
  body: API.PatchPageDataRequest[],
  options?: { [key: string]: any },
) {
  const { pageId: param0, ...queryParams } = params;
  return request<API.BaseTotalResponse>(`/api/v1/pages/${param0}/data`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** Get the specified data of a page. GET /api/v1/pages/${param0}/data/${param1} */
export async function getPageData(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getPageDataParams,
  options?: { [key: string]: any },
) {
  const { pageId: param0, id: param1, ...queryParams } = params;
  return request<API.GetPageDataResponse>(`/api/v1/pages/${param0}/data/${param1}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** Update data of a page. (full). PUT /api/v1/pages/${param0}/data/${param1} */
export async function updatePageData(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.updatePageDataParams,
  body: API.UpdatePageDataRequest,
  options?: { [key: string]: any },
) {
  const { pageId: param0, id: param1, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/pages/${param0}/data/${param1}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** Delete data of a page. DELETE /api/v1/pages/${param0}/data/${param1} */
export async function deletePageData(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deletePageDataParams,
  options?: { [key: string]: any },
) {
  const { pageId: param0, id: param1, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/pages/${param0}/data/${param1}`, {
    method: 'DELETE',
    params: { ...queryParams },
    ...(options || {}),
  });
}
