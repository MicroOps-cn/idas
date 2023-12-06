// @ts-ignore

/* eslint-disable */
import { request } from '@/utils/request';

/** Obtain Security Configuration. GET /api/v1/config/security */
export async function getSecurityConfig(options?: { [key: string]: any }) {
  return request<API.GetSecurityConfigResponse>('/api/v1/config/security', {
    method: 'GET',
    ...(options || {}),
  });
}

/** Update Security Configuration (Incremental). PATCH /api/v1/config/security */
export async function patchSecurityConfig(
  body: API.PatchSecurityConfigRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponse>('/api/v1/config/security', {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
