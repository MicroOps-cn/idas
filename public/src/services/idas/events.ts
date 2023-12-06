// @ts-ignore

/* eslint-disable */
import { request } from '@/utils/request';

/** Get events. GET /api/v1/events */
export async function getEvents(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getEventsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetEventsResponse>('/api/v1/events', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Get event logs. GET /api/v1/events/logs */
export async function getEventLogs(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getEventLogsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetEventLogsResponse>('/api/v1/events/logs', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
