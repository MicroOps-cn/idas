// @ts-ignore

/* eslint-disable */
import { request } from '@/utils/request';

/** Get permission list. GET /api/v1/permissions */
export async function getPermissions(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getPermissionsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetPermissionsResponse>('/api/v1/permissions', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
