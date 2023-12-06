import { loginPath, forgotPasswordPath } from './env';

const routes: Route[] = [
  {
    path: '/account',
    layout: false,
    routes: [
      {
        path: loginPath,
        component: './Account/Login',
        access: 'canAnonymous',
        // // 不展示顶栏
        // headerRender: false,
        // // 不展示页脚
        // footerRender: false,
        // // 不展示菜单
        // menuRender: false,
        // // 不展示菜单顶栏
        // menuHeaderRender: false,
      },
      {
        path: forgotPasswordPath,
        component: './Account/ForgotPassword',
        access: 'canAnonymous',
      },
      {
        path: '/account/resetPassword',
        component: './Account/ResetPassword',
        access: 'canAnonymous',
      },
      {
        path: '/account/activateAccount',
        component: './Account/ActivateAccount',
        access: 'canAnonymous',
      },
      {
        path: '/account/settings',
        component: './Account/Setting',
      },
      {
        path: '/account/events',
        component: './Account/Events',
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/welcome',
    layout: false,
    access: 'canUser',
    component: './Welcome',
  },
  {
    path: '/home',
    redirect: '/welcome',
    access: 'canUser',
    name: 'home',
    icon: 'HomeOutlined',
  },
  // {
  //   path: '/dashboard',
  //   name: 'dashboard',
  //   icon: 'dashboard',
  //   access: 'canViewer',
  //   component: './Dashboard',
  // },
  {
    name: 'apps',
    icon: 'AppstoreOutlined',
    path: '/apps',
    exact: false,
    routes: [
      {
        exact: true,
        path: '/apps',
        access: 'canViewer',
        component: './Apps/index',
      },
      {
        exact: true,
        access: 'canViewer',
        path: '/apps/create',
        component: './Apps/[aid].edit',
      },
      {
        exact: true,
        access: 'canViewer',
        path: '/apps/:aid',
        component: './Apps/[aid]',
      },
      {
        exact: true,
        access: 'canEditor',
        path: '/apps/:aid/edit',
        component: './Apps/[aid].edit',
      },
    ],
  },
  {
    name: 'roles',
    icon: 'table',
    path: '/roles',
    access: 'canViewer',
    component: './Roles',
  },
  {
    name: 'users',
    icon: 'UserOutlined',
    path: '/users',
    access: 'canViewer',
    component: './User',
  },
  {
    name: 'events',
    path: '/events',
    component: './Events',
    icon: 'SecurityScanOutlined',
    access: 'canAdmin',
  },
  {
    name: 'settings',
    path: '/settings',
    icon: 'setting',
    access: 'canAdmin',
    component: './Setting',
  },
  // {
  //   name: 'pages',
  //   icon: 'UserOutlined',
  //   path: '/pages',
  //   access: 'canAdmin',
  //   routes: [
  //     {
  //       exact: true,
  //       path: '/pages',
  //       access: 'canAdmin',
  //       component: './Pages/index',
  //     },
  //     {
  //       exact: true,
  //       path: '/pages/:pageId',
  //       access: 'canAdmin',
  //       component: './Pages/[pageId]',
  //     },
  //   ],
  // },
  // {
  //   path: '/page',
  //   routes: [
  //     {
  //       exact: true,
  //       path: '/page/:pageId/:id',
  //       access: 'canAdmin',
  //       component: './Page/[pageId][id]',
  //     },
  //     {
  //       exact: true,
  //       path: '/page/:pageId',
  //       access: 'canAdmin',
  //       component: './Page/[pageId]',
  //     },
  //   ],
  // },
  {
    path: '/403',
    layout: false,
    access: 'canAnonymous',
    component: './403',
  },
  {
    path: '/warning',
    layout: false,
    access: 'canAnonymous',
    component: './Warning',
  },
  {
    path: '/',
    redirect: '/welcome',
  },
  {
    component: './404',
  },
];

export interface Route {
  path?: string;
  component?: string | (() => any);
  wrappers?: string[];
  redirect?: string;
  exact?: boolean;
  routes?: Route[];
  access?: 'canViewer' | 'canUser' | 'canAdmin' | 'canEditor' | 'canAnonymous' | 'forbidden';
  [k: string]: any;
}

export default routes;
