// @ts-ignore

/* eslint-disable */
import { request } from '@/utils/request';

/** Well known openid configuration. GET /api/v1/oauth/.well-known/openid-configuration */
export async function openidConfiguration(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.openidConfigurationParams,
  options?: { [key: string]: any },
) {
  return request<API.OIDCWellKnownResponse>('/api/v1/oauth/.well-known/openid-configuration', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Application authorization. GET /api/v1/oauth/authorize */
export async function oAuthAuthorize(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.oAuthAuthorizeParams,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponse>('/api/v1/oauth/authorize', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Application authorization. POST /api/v1/oauth/authorize */
export async function oAuthAuthorize2(
  body: API.OAuthAuthorizeRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponse>('/api/v1/oauth/authorize', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** Well known openid configuration. GET /api/v1/oauth/jwks */
export async function openidConfiguration2(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.openidConfigurationParams,
  options?: { [key: string]: any },
) {
  return request<API.OAuthJWKSResponse>('/api/v1/oauth/jwks', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Get token. POST /api/v1/oauth/token */
export async function oAuthTokens(body: API.OAuthTokenRequest, options?: { [key: string]: any }) {
  const formData = new FormData();

  Object.keys(body).forEach((ele) => {
    const item = (body as any)[ele];

    if (item !== undefined && item !== null) {
      if (typeof item === 'object' && !(item instanceof File)) {
        if (item instanceof Array) {
          item.forEach((f) => formData.append(ele, f || ''));
        } else {
          formData.append(ele, JSON.stringify(item));
        }
      } else {
        formData.append(ele, item);
      }
    }
  });

  return request<API.OAuthTokenResponse>('/api/v1/oauth/token', {
    method: 'POST',
    data: formData,
    ...(options || {}),
  });
}

/** Get user info. GET /api/v1/oauth/userinfo */
export async function oAuthUserInfo(options?: { [key: string]: any }) {
  return request<API.BaseResponse>('/api/v1/oauth/userinfo', {
    method: 'GET',
    ...(options || {}),
  });
}
