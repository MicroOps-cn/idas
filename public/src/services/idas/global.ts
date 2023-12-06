// @ts-ignore

/* eslint-disable */
import { request } from '@/utils/request';

/** Get global config. GET /api/v1/global/config */
export async function getGlobalConfig(options?: { [key: string]: any }) {
  return request<API.GlobalConfigResponse>('/api/v1/global/config', {
    method: 'GET',
    ...(options || {}),
  });
}
