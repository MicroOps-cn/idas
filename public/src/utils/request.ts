// decoupling with antd UI library, you can using `alias` modify the ui methods
import { message, notification } from 'antd';

import { request as umiRequest, history } from '@umijs/max';
import type {
  AxiosError,
  AxiosRequestConfig,
  Request,
  RequestError as UmiRequestError,
  RequestOptions as IRequestOptions,
} from '@umijs/max';

import type { IntlContext } from './intl';

// @ts-ignore
// import { request as umiRequest,RequestOptions as IRequestOptions,Request } from 'umi';
// import type {
//   Context,
//   RequestOptionsWithResponse,
//   RequestOptionsInit,
//   RequestOptionsWithoutResponse,
// } from 'umi-request';

// //@ts-ignore
// import { message, notification } from 'antd';

// import type { IntlContext } from './intl';
// import { AxiosResponse } from '@umijs/max';

export enum ErrorShowType {
  SILENT = 0,
  WARN_MESSAGE = 1,
  ERROR_MESSAGE = 2,
  NOTIFICATION = 4,
  REDIRECT = 9,
}
export type RequestError = UmiRequestError & {
  data?: any;
  info?: ResponseStructure;
  handled?: boolean;
  config?: AxiosRequestConfig<any> & { intl?: IntlContext };
};
export interface ResponseStructure {
  success: boolean;
  data?: any;
  errorCode?: string;
  errorMessage?: string;
  showType?: ErrorShowType;
  traceId?: string;
  [key: string]: any;
}

const errorAdaptor = (resData: any) => {
  if (typeof resData === 'string') {
    return {
      errorMessage: resData,
    };
  }
  return resData;
};

const DEFAULT_ERROR_PAGE = '/warning';
export const errorHandler = (error: RequestError, opts: IRequestOptions) => {
  // @ts-ignore
  if (opts.ignoreError) {
    return
  }
  if (opts.skipErrorHandler) {
    throw error;
  }
  let errorInfo: ResponseStructure | undefined;
  if (error.name === 'AxiosError' && (error as AxiosError).response?.data) {
    errorInfo = errorAdaptor((error as AxiosError).response?.data);
    error.message = errorInfo?.errorMessage || error.message;
    error.info = errorInfo;
  } else if (error.name === 'AxiosError') {
    errorInfo = error.info;
  }
  if (errorInfo) {
    const errorCode = errorInfo?.errorCode;
    const intl = error?.config?.intl;
    const errorMessage =
      intl && errorCode
        ? intl.t(`${errorCode}.errorMessage`, errorInfo?.errorMessage)
        : errorInfo?.errorMessage;
    const errorPage = DEFAULT_ERROR_PAGE;
    switch (errorInfo?.showType) {
      case ErrorShowType.SILENT:
        // do nothing
        break;
      case ErrorShowType.WARN_MESSAGE:
        message.warning(errorMessage);
        break;
      case ErrorShowType.ERROR_MESSAGE:
        message.error(errorMessage);
        break;
      case ErrorShowType.NOTIFICATION:
        notification.open({
          description: errorMessage,
          message: errorCode,
        });
        break;
      case ErrorShowType.REDIRECT:
        // @ts-ignore
        history.push({
          pathname: errorPage,
          query: { errorCode, errorMessage },
        });
        // redirect to error page
        break;
      default:
        message.error(errorMessage);
        break;
    }
  } else {
    message.error(error.message || 'Request error, please retry.');
  }
  error.handled = true;
  throw error;
};

export type RequestOptions = IRequestOptions & {
  getResponse?: boolean;
  skipErrorHandler?: boolean;
  ignoreError?: boolean;
};

declare const apiPath: string;
declare const basePath: string;

export const getApiPath = (api: string): string => {
  if (apiPath) {
    if (apiPath.endsWith('/') || api.startsWith('/')) {
      if (apiPath.endsWith('/') && api.startsWith('/')) {
        return apiPath.slice(0, -1) + api;
      }
      return apiPath + api;
    }
    return apiPath + '/' + api;
  }
  return api.startsWith('/') ? api : '/' + api;
};

export const getLocation = () => {
  const location = history.location;
  const { pathname } = location;
  if (basePath && pathname.startsWith(basePath)) {
    const p = pathname.slice(basePath.length);
    return { ...location, pathname: p.startsWith('/') ? p : `/${p}` };
  }
  return location;
};

// export const request: Request = (url: string, options?: RequestOptions) => umiRequest(getApiPath(apiPath), options as any);

export const request: Request = (url: string, options?: RequestOptions) =>
  umiRequest(url, options as any);
export default request;

declare const publicPath: string;

export const getPublicPath = (subPath: string): string => {
  if (publicPath) {
    if (publicPath.endsWith('/') || subPath.startsWith('/')) {
      if (publicPath.endsWith('/') && subPath.startsWith('/')) {
        return publicPath.slice(0, -1) + subPath;
      }
      return publicPath + subPath;
    }
    return publicPath + '/' + subPath;
  }
  return subPath.startsWith('/') ? subPath : '/' + subPath;
};
