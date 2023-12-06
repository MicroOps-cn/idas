import Footer from '@/components/Footer';
import { getActions } from '@/components/RightContent';
import type { ResponseStructure } from '@/utils/request';
import { errorHandler, getApiPath, getLocation, getPublicPath } from '@/utils/request';
import { BookOutlined, LinkOutlined } from '@ant-design/icons';
import type { MenuDataItem, Settings as LayoutSettings } from '@ant-design/pro-components';
import { SettingDrawer } from '@ant-design/pro-components';
import type { AxiosResponse, RequestConfig, RunTimeLayoutConfig } from '@umijs/max';
import { history, Link } from '@umijs/max';

import defaultSettings from '../config/defaultSettings';
import { loginPath } from '../config/env';
import type { Route } from '../config/routes';
import globalRoutes from '../config/routes';
import ForbiddenPage from './pages/403';
import NoFoundPage from './pages/404';
import { getGlobalConfig } from './services/idas/global';
import { currentUser as queryCurrentUser } from './services/idas/user';

const isDev = process.env.NODE_ENV === 'development';

function getRouteAccess(path: string, routes: Route[]): { ok: boolean; access?: string } {
  for (const route of routes) {
    if (route.path === path || (route.path ?? '') + '/' === path) {
      return { ok: true, access: route.access };
    }
    if (route.routes) {
      const { ok, access } = getRouteAccess(path, route.routes);
      if (ok) {
        return { ok, access: access ? access : route.access };
      }
    }
  }
  return { ok: false };
}

/**
 * @see  https://umijs.org/zh-CN/plugins/plugin-initial-state
 * */
export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  globalConfig?: API.GlobalConfig;
  currentUser?: API.UserInfo & { fetchTime: number };
  loading?: boolean;
  fetchUserInfo?: () => Promise<(API.UserInfo & { fetchTime: number }) | undefined>;
  dynamicMenu?: MenuDataItem[];
}> {
  const { data: globalConfig } = await getGlobalConfig();
  const { pathname } = getLocation();
  const settings = { ...defaultSettings };
  if (globalConfig?.logo) {
    settings.logo = globalConfig.logo;
  }
  if (!settings.logo) {
    settings.logo = getPublicPath('logo.svg');
  }
  if (globalConfig?.title) {
    settings.title = globalConfig.title;
  }
  const fetchUserInfo: () => Promise<
    (API.UserInfo & { fetchTime: number }) | undefined
  > = async () => {
    try {
      const msg = await queryCurrentUser({ skipErrorHandler: pathname === '/' });
      return msg.data ? { ...msg.data, fetchTime: Number(new Date()) / 1000 } : undefined;
    } catch (error) {
      history.push(loginPath);
    }
    return undefined;
  };
  const { access } = getRouteAccess(pathname, globalRoutes);
  if (access !== 'canAnonymous') {
    const currentUser = await fetchUserInfo();
    return {
      fetchUserInfo,
      currentUser,
      settings,
      globalConfig,
    };
  }
  return {
    fetchUserInfo,
    settings,
    globalConfig,
  };
}

// ProLayout 支持的api https://procomponents.ant.design/components/layout
export const layout: RunTimeLayoutConfig = ({ initialState, setInitialState }) => {
  const { pathname } = getLocation();
  return {
    actionsRender: () => getActions(initialState?.currentUser?.role),
    disableContentMargin: false,
    // waterMarkProps: { // 水印
    //   content: initialState?.currentUser?.username,
    // },
    footerRender: () => <Footer />,
    onPageChange: () => {
      // 如果没有登录，重定向到 login
      const { access } = getRouteAccess(pathname, globalRoutes);
      if (!initialState?.currentUser && access !== 'canAnonymous') {
        history.push(loginPath);
      }
    },
    links: isDev
      ? [
          <Link key="openapi" to="/umi/plugin/openapi" target="_blank">
            <LinkOutlined />
            <span>OpenAPI 文档</span>
          </Link>,
          <Link key="docs" to="/~docs">
            <BookOutlined />
            <span>业务组件文档</span>
          </Link>,
        ]
      : [],
    menuDataRender(menuData) {
      if (initialState?.currentUser) {
        if (initialState.dynamicMenu === undefined) {
          // MenuRender.getInstance().render((dynamicMenu) => {
          //   if (dynamicMenu.length > 0) {
          //     setInitialState({
          //       ...initialState,
          //       dynamicMenu,
          //     });
          //   }
          // });
        }
      }
      return [...menuData, ...(initialState?.dynamicMenu ?? [])];
    },
    menuHeaderRender: undefined,
    // 自定义 403 页面
    unAccessible: <ForbiddenPage />,
    noFound: <NoFoundPage />,
    // 增加一个 loading 的状态
    childrenRender: (children, props) => {
      // if (initialState?.loading) return <PageLoading />;
      return (
        <>
          {children}
          {isDev && !props.location?.pathname?.includes('/login') && (
            <SettingDrawer
              enableDarkTheme
              settings={initialState?.settings}
              onSettingChange={(settings) => {
                setInitialState((preInitialState) => ({
                  ...preInitialState,
                  settings: {
                    ...initialState?.settings,
                    ...settings,
                  },
                }));
              }}
            />
          )}
        </>
      );
    },
    ...initialState?.settings,
  };
};

/**
 * @name request 配置，可以配置错误处理
 * 它基于 axios 和 ahooks 的 useRequest 提供了一套统一的网络请求和错误处理方案。
 * @doc https://umijs.org/docs/max/request#配置
 */
export const request: RequestConfig = {
  requestInterceptors: [
    (url: any, options: any) => {
      return { url: getApiPath(url), options };
    },
  ],
  responseInterceptors: [
    (response) => {
      const { config, data } = response as AxiosResponse<ResponseStructure>;

      const { success, errorMessage, errorCode } = data;
      if (!success) {
        const error: any = new Error(errorMessage);
        error.name = 'BizError';
        error.info = data;
        error.code = errorCode;
        error.config = config;
        error.response = response;
        error.request = response.request;
        throw error; // 抛出自制的错误
      }
      return response;
    },
  ],
  errorConfig: {
    errorHandler: errorHandler,
  },
};
