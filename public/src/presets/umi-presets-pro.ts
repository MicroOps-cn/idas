import type { IApi } from '@umijs/max';

export default (api: IApi) => {
  api.onStart(() => {
    console.log('😄 Hello PRO');
  });
  const plugins = [
    require.resolve('umi-presets-pro/dist/features/proconfig'),
    require.resolve('umi-presets-pro/dist/features/maxtabs'),
    require.resolve('../plugins/openapi'),
    require.resolve('@alita/plugins/dist/keepalive'),
    require.resolve('@alita/plugins/dist/tabs-layout'),
    require.resolve('@umijs/request-record'),
  ];
  return {
    plugins,
  };
};
