import { message } from 'antd';
import React from 'react';
import { useIntl, history, useModel } from 'umi';

import defaultSettings from '@/../config/defaultSettings';
import Footer from '@/components/Footer';
import List from '@/components/List';
import { getActions } from '@/components/RightContent';
import { currentUserApps } from '@/services/idas/user';
import { IntlContext } from '@/utils/intl';
import { getPublicPath } from '@/utils/request';
import { ProLayout } from '@ant-design/pro-components';

import styles from './Welcome.less';

const Welcome: React.FC = () => {
  const intl = new IntlContext('pages.welcome', useIntl());
  const {
    initialState = {
      settings: { navTheme: 'light' },
      globalConfig: null,
      currentUser: { role: undefined },
    },
  } = useModel('@@initialState');
  // @ts-expect-error
  const { navTheme = 'dark' } = initialState.settings;
  const globalConfig = initialState?.globalConfig ?? null;
  return (
    <ProLayout
      logo={globalConfig?.logo ?? getPublicPath('logo.svg')}
      title={globalConfig?.title ?? defaultSettings.title}
      onMenuHeaderClick={() => {
        history.push('/');
      }}
      layout="top"
      navTheme={navTheme}
      actionsRender={() => getActions(initialState?.currentUser?.role)}
    >
      <List<API.AppInfo>
        cardProps={{
          className: styles.AppList,
        }}
        intl={intl}
        request={{
          list: currentUserApps,
        }}
        onClick={(item) => {
          if (item.url) {
            const w = window.open('about:blank');
            if (w) w.location.href = item.url;
          } else {
            message.warn(
              'The URL for this application is not configured. Please contact the administrator.',
            );
          }
        }}
      />
      <Footer className={styles.Footer} />
    </ProLayout>
  );
};

export default Welcome;
