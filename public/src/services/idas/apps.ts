// @ts-ignore

/* eslint-disable */
import { request } from '@/utils/request';

/** Get the application list. GET /api/v1/apps */
export async function getApps(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAppsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetAppsResponse>('/api/v1/apps', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Create an application. POST /api/v1/apps */
export async function createApp(body: API.CreateAppRequest, options?: { [key: string]: any }) {
  return request<any>('/api/v1/apps', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Batch delete applications. DELETE /api/v1/apps */
export async function deleteApps(body: API.DeleteAppRequest[], options?: { [key: string]: any }) {
  return request<API.BaseTotalResponse>('/api/v1/apps', {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Batch update of application information (incremental). PATCH /api/v1/apps */
export async function patchApps(body: API.PatchAppRequest[], options?: { [key: string]: any }) {
  return request<API.BaseTotalResponse>('/api/v1/apps', {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Get Application info. GET /api/v1/apps/${param0} */
export async function getAppInfo(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAppInfoParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.GetAppResponse>(`/api/v1/apps/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新应用信息（全量） PUT /api/v1/apps/${param0} */
export async function updateApp(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.updateAppParams,
  body: API.UpdateAppRequest,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/apps/${param0}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** Delete app. DELETE /api/v1/apps/${param0} */
export async function deleteApp(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteAppParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/apps/${param0}`, {
    method: 'DELETE',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** Update application information (full). PATCH /api/v1/apps/${param0} */
export async function patchApp(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.patchAppParams,
  body: API.PatchAppRequest,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/apps/${param0}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** Get a app key-pairs. GET /api/v1/apps/${param0}/key */
export async function getAppKeys(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAppKeysParams,
  options?: { [key: string]: any },
) {
  const { appId: param0, ...queryParams } = params;
  return request<API.GetAppKeysResponse>(`/api/v1/apps/${param0}/key`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** Create a app key pair. POST /api/v1/apps/${param0}/key */
export async function createAppKey(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.createAppKeyParams,
  body: API.CreateAppKeyRequest,
  options?: { [key: string]: any },
) {
  const { appId: param0, ...queryParams } = params;
  return request<API.CreateAppKeyResponse>(`/api/v1/apps/${param0}/key`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** Delete a app key pairs. DELETE /api/v1/apps/${param0}/key */
export async function deleteAppKeys(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteAppKeysParams,
  body: API.DeleteAppKeysRequest,
  options?: { [key: string]: any },
) {
  const { appId: param0, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/apps/${param0}/key`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** Get a app Icons. GET /api/v1/apps/icons */
export async function getAppIcons(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAppIconsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetAppIconsResponse>('/api/v1/apps/icons', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
