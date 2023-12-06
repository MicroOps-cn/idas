// @ts-ignore

/* eslint-disable */
import { request } from '@/utils/request';

/** Upload file POST /api/v1/files */
export async function uploadFile(options?: { [key: string]: any }) {
  return request<API.FileUploadResponse>('/api/v1/files', {
    method: 'POST',
    ...(options || {}),
  });
}

/** Download/View File GET /api/v1/files/${param0} */
export async function downloadFile(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.downloadFileParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/files/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}
