// @ts-ignore

/* eslint-disable */
import { request } from '@/utils/request';

/** Get current login user information. GET /api/v1/user */
export async function currentUser(options?: { [key: string]: any }) {
  return request<API.GetUserResponse>('/api/v1/user', {
    method: 'GET',
    ...(options || {}),
  });
}

/** Update current login user information (full). PUT /api/v1/user */
export async function updateCurrentUser(
  body: API.UpdateUserRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponse>('/api/v1/user', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Update current login user information (increment). PATCH /api/v1/user */
export async function patchCurrentUser(
  body: API.PatchCurrentUserRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponse>('/api/v1/user', {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Activate the user. POST /api/v1/user/activateAccount */
export async function activateAccount(
  body: API.ActivateAccountRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponse>('/api/v1/user/activateAccount', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Get current login user's apps. GET /api/v1/user/apps */
export async function currentUserApps(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.currentUserAppsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetAppsResponse>('/api/v1/user/apps', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Get current login user's event logs. GET /api/v1/user/eventLogs */
export async function currentUserEventLogs(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.currentUserEventLogsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetCurrentUserEventLogsResponse>('/api/v1/user/eventLogs', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Get current login user's events. GET /api/v1/user/events */
export async function currentUserEvents(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.currentUserEventsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetCurrentUserEventsResponse>('/api/v1/user/events', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Forgot the user password. POST /api/v1/user/forgotPassword */
export async function forgotPassword(
  body: API.ForgotUserPasswordRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponse>('/api/v1/user/forgotPassword', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** User login. POST /api/v1/user/login */
export async function userLogin(body: API.UserLoginRequest, options?: { [key: string]: any }) {
  return request<API.UserLoginResponse>('/api/v1/user/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** User logout. POST /api/v1/user/logout */
export async function userLogout(options?: { [key: string]: any }) {
  return request<API.BaseResponse>('/api/v1/user/logout', {
    method: 'POST',
    ...(options || {}),
  });
}

/** OAuth login. GET /api/v1/user/oauth/${param0} */
export async function userOAuthLogin(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.userOAuthLoginParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<any>(`/api/v1/user/oauth/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** Reset the user password. POST /api/v1/user/resetPassword */
export async function resetPassword(
  body: API.ResetUserPasswordRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponse>('/api/v1/user/resetPassword', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Send login code. POST /api/v1/user/sendLoginCaptcha */
export async function sendLoginCaptcha(
  body: API.SendLoginCaptchaRequest,
  options?: { [key: string]: any },
) {
  return request<API.SendLoginCaptchaResponse>('/api/v1/user/sendLoginCaptcha', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Get current user session list. GET /api/v1/user/sessions */
export async function getCurrentUserSessions(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getCurrentUserSessionsParams,
  options?: { [key: string]: any },
) {
  return request<API.GetSessionsResponse>('/api/v1/user/sessions', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Delete current user a session. DELETE /api/v1/user/sessions/${param0} */
export async function deleteCurrentUserSession(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteCurrentUserSessionParams,
  options?: { [key: string]: any },
) {
  const { id: param0, ...queryParams } = params;
  return request<API.BaseResponse>(`/api/v1/user/sessions/${param0}`, {
    method: 'DELETE',
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** binding TOTP Secret POST /api/v1/user/totp */
export async function bindingTotp(body: API.CreateTOTPRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponse>('/api/v1/user/totp', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** get TOTP Secret GET /api/v1/user/totp/secret */
export async function getTotpSecret(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getTOTPSecretParams,
  options?: { [key: string]: any },
) {
  return request<API.CreateTOTPSecretResponse>('/api/v1/user/totp/secret', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
