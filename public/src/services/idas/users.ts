// @ts-ignore

/* eslint-disable */
import { request } from '@/utils/request';

/** Get user list. GET /api/v1/users */
export async function getUsers(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getUsersParams,
  options?: { [key: string]: any },
) {
  return request<API.GetUsersResponse>('/api/v1/users', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Create a user. POST /api/v1/users */
export async function createUser(body: API.CreateUserRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponse>('/api/v1/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Delete users in batch. DELETE /api/v1/users */
export async function deleteUsers(body: API.DeleteUserRequest[], options?: { [key: string]: any }) {
  return request<API.BaseTotalResponse>('/api/v1/users', {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Batch update user information(Incremental). PATCH /api/v1/users */
export async function patchUsers(body: API.PatchUserRequest[], options?: { [key: string]: any }) {
  return request<API.BaseTotalResponse>('/api/v1/users', {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Get user information. GET /api/v1/users/${param0} */
export async function getUserInfo(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getUserInfoParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.GetUserResponse>(`/api/v1/users/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** Update user information(full). PUT /api/v1/users/${param0} */
export async function updateUser(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.updateUserParams,
  body: API.UpdateUserRequest,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.UpdateUserRequest>(`/api/v1/users/${param0}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** Delete user. DELETE /api/v1/users/${param0} */
export async function deleteUser(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteUserParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/users/${param0}`, {
    method: 'DELETE',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** Update user information(Incremental). PATCH /api/v1/users/${param0} */
export async function patchUser(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.patchUserParams,
  body: API.PatchUserRequest,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.PatchUserResponse>(`/api/v1/users/${param0}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** Send account activation email. POST /api/v1/users/sendActivateMail */
export async function sendActivateMail(
  body: API.SendActivationMailRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponse>('/api/v1/users/sendActivateMail', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
