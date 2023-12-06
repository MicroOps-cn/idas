// @ts-ignore

/* eslint-disable */
import { request } from '@/utils/request';

/** Get session list. GET /api/v1/sessions */
export async function getSessions(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getSessionsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetSessionsResponse>('/api/v1/sessions', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Expire a session. DELETE /api/v1/sessions/${param0} */
export async function deleteSession(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteSessionParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/sessions/${param0}`, {
    method: 'DELETE',
    params: { ...queryParams },
    ...(options || {}),
  });
}
